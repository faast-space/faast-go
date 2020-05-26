package v1

import (
	"encoding/json"

	"github.com/opentracing/opentracing-go"
)

type (
	iContext struct {
		event  iEvent
		trace  opentracing.SpanContext
		answer map[string]interface{}
	}
	iEvent struct {
		source     Source
		parameters map[string]interface{}
	}
	iParameterValue struct {
		v interface{}
	}
)

func (c *iContext) Event() Event {
	return &c.event
}

func (c *iContext) Answer(m map[string]interface{}) error {
	c.answer = m
	return nil
}

func (c *iContext) Trace() opentracing.SpanContext {
	return c.trace
}

func (e *iEvent) From() Source {
	return e.source
}

func (e *iEvent) Parameter(name string) ParameterValue {
	return iParameterValue{
		v: e.parameters[name],
	}
}

func (p iParameterValue) String() string {
	s, _ := p.v.(string)
	return s
}

func (p iParameterValue) Int64() int64 {
	i, _ := p.v.(int64)
	return i
}

func (p iParameterValue) Float64() float64 {
	f, _ := p.v.(float64)
	return f
}

func (p iParameterValue) Bool() bool {
	ok, _ := p.v.(bool)
	return ok
}

func (p iParameterValue) JSONUnmarshal(v interface{}) error {
	b, err := json.Marshal(p.v)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, v)
}
