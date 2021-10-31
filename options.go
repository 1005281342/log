package log

import (
	"context"
	"go.opentelemetry.io/otel/trace"
)

// Option ...
type Option struct {
	Address  string
	AppName  string
	FuncName string
	TraceID  string
}

// OptionFunc ...
type OptionFunc func(opt *Option)

// WithAddress with address
func WithAddress(addr string) func(opt *Option) {
	return func(opt *Option) {
		opt.Address = addr
	}
}

// WithFuncName with func name
func WithFuncName(name string) func(opt *Option) {
	return func(opt *Option) {
		opt.FuncName = name
	}
}

// WithAppName with app name
func WithAppName(name string) func(opt *Option) {
	return func(opt *Option) {
		opt.AppName = name
	}
}

// WithTraceID with traceID
func WithTraceID(traceID string) func(opt *Option) {
	return func(opt *Option) {
		opt.TraceID = traceID
	}
}

func traceIdFromContext(ctx context.Context) string {
	spanCtx := trace.SpanContextFromContext(ctx)
	if spanCtx.HasTraceID() {
		return spanCtx.TraceID().String()
	}

	return ""
}
