package v1

// #cgo CFLAGS: -I"C:/Program Files/foundationdb/include"
// #cgo LDFLAGS: -L"C:/Program Files/foundationdb/bin" -lfdb_c
// #define FDB_API_VERSION 620
// #include <foundationdb/fdb_c.h>

import (
	"C"
	"context"
	"net"

	"github.com/faast-space/faast-go/v1/function"
	"github.com/faast-space/faast-go/v1/log"
	jaeger "github.com/uber/jaeger-client-go"
	"google.golang.org/grpc"
)
import (
	"encoding/json"
	"os"
	"os/signal"
	"syscall"

	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc/metadata"
)

type (
	// Handler is a high level function definition
	Handler func(Context) error

	// Source tell who trigger the function
	Source string

	// Event contains all parameters for execution
	Event interface {
		// Who made the call
		From() Source
		// Get parameter value
		Parameter(string) ParameterValue
	}

	// Context is the object pass to function
	Context interface {
		Event() Event
		// Get trace identifier
		Trace() opentracing.SpanContext
		Answer(map[string]interface{}) error
	}

	// ParameterValue is an helpful class to parse your parameters
	ParameterValue interface {
		String() string
		Int64() int64
		Float64() float64
		Bool() bool
		JSONUnmarshal(v interface{}) error
	}

	rpcServer struct {
		function.UnimplementedFunctionServer
		fn Handler
	}
)

const (
	// SourceHTTP is an HTTP call
	SourceHTTP Source = "HTTP"
	// SourceStream is an event reaction
	SourceStream Source = "STREAM"
	// SourceCRON is a scheduled execution
	SourceCRON Source = "CRON"
)

func (s *rpcServer) Execute(ctx context.Context, in *function.Request) (*function.Response, error) {
	var err error
	log.Info("Received: %v", in.GetKind())

	var span jaeger.SpanContext

	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		log.Info("%+v", md)
		if values := md.Get("trace"); len(values) > 0 {
			span, err = jaeger.ContextFromString(values[0])
			if err != nil {
				log.WithError(err).Warn("no trace context provided")
			}
		}
	} else {
		log.Warn("no metadatas")
	}

	parameters := map[string]interface{}{}
	if in.GetParameters() != nil {
		err = json.Unmarshal(in.GetParameters(), &parameters)
		if err != nil {
			log.WithError(err).Error("cannot parse payload")
		}
	}

	ev := iEvent{
		source:     Source(in.GetKind().String()),
		parameters: parameters,
	}

	fnCtx := &iContext{
		event:  ev,
		trace:  span,
		answer: map[string]interface{}{},
	}

	err = s.fn(fnCtx)

	errStr := ""
	if err != nil {
		errStr = err.Error()
	}

	if fnCtx.answer == nil {
		fnCtx.answer = map[string]interface{}{}
	}
	b, err := json.Marshal(fnCtx.answer)
	if err != nil {
		log.WithError(err).Error("cannot serialize response")
	}

	return &function.Response{
		Error:      errStr,
		Parameters: b,
	}, nil
}

// Start listening for events
func Start(fn Handler) {
	lis, err := net.Listen("tcp", "0.0.0.0:8080")
	if err != nil {
		log.WithError(err).Fatal("failed to listen")
	}

	s := grpc.NewServer()
	function.RegisterFunctionServer(s, &rpcServer{
		fn: fn,
	})

	stop := make(chan os.Signal, 2)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := s.Serve(lis); err != nil {
			log.WithError(err).Fatal("failed to serve")
		}
	}()

	sig := <-stop
	log.Warn("received %v", sig)

	s.GracefulStop()
}
