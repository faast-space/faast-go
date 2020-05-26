package test

import (
	faast "github.com/faast-space/faast-go/v1"
	"github.com/uber/jaeger-client-go"
)

type (
	Context struct {
		event Event
		trace jaeger.SpanContext
	}

	Event struct {
		from faast.Source
	}
)

func NewContext() faast.Context {
	return &Context{
		event: Event{},
		trace: jaeger.SpanContext{},
	}
}

func (c *Context) Event() faast.Event {
	return &c.event
}

func (c *Context) JSON(object interface{}) error {
	return nil
}

func (c *Context) Trace() jaeger.SpanContext {
	return c.trace
}

func (e *Event) From() faast.Source {
	return e.from
}
