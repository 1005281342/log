package log

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/bshuster-repo/logrus-logstash-hook"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/tal-tech/go-zero/core/logx"
)

const (
	TraceIDKey = "traceID"
)

// NewXELKLoggerWithContext new go-zero elk logger
func NewXELKLoggerWithContext(ctx context.Context, optFuncList ...OptionFunc) logx.Logger {
	var (
		logger         logx.Logger
		loggerELK, err = WithContext(ctx, optFuncList...)
	)
	if err != nil {
		logger = logx.WithContext(ctx)
		logger.Errorf("new ELKLogger WithContext err: %+v", err)
	} else {
		logger = &XELKLogger{ELKLogger: loggerELK}
	}
	return logger
}

type XELKLogger struct {
	*ELKLogger
}

func (e *XELKLogger) WithDuration(d time.Duration) logx.Logger {
	e.ELKLogger.WithDuration(d)
	return e
}

type ELKLogger struct {
	logrus.FieldLogger
	Duration string
	*Option
}

// WithContext new elk logger
func WithContext(ctx context.Context, optFuncList ...OptionFunc) (*ELKLogger, error) {
	var (
		logger = logrus.New()
		conn   net.Conn
		err    error
		elk    = &ELKLogger{Option: &Option{Address: "127.0.0.1:5000"}}
	)

	for _, optFunc := range optFuncList {
		optFunc(elk.Option)
	}
	var fields = logrus.Fields{
		"AppName":  elk.AppName,
		"FuncName": elk.FuncName,
	}

	if conn, err = net.Dial("tcp", elk.Address); err != nil {
		return nil, err
	}
	var hook = logrustash.New(conn, logrustash.DefaultFormatter(fields))
	logger.Hooks.Add(hook)

	var entry = logger.WithFields(logrus.Fields{
		TraceIDKey: getTraceID(ctx),
	})
	elk.FieldLogger = entry
	return elk, nil
}

func WithContextAndAddress(ctx context.Context, addr string) (*ELKLogger, error) {
	return WithContext(ctx, WithAddress(addr))
}

func SetTraceID(ctx context.Context, traceID string) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithValue(ctx, TraceIDKey, traceID)
}

func getTraceID(ctx context.Context) string {
	if ctx == nil {
		return uuid.NewString()
	}

	var v, ok = (ctx.Value(TraceIDKey)).(string)
	if !ok || v == "" {
		return uuid.NewString()
	}
	return v
}

func (e *ELKLogger) Errorv(v interface{}) {
	e.FieldLogger.Error(v)
}

func (e *ELKLogger) Infov(v interface{}) {
	e.FieldLogger.Info(v)
}

func (e *ELKLogger) Slow(args ...interface{}) {
	e.FieldLogger.Info(args...)
}

func (e *ELKLogger) Slowf(s string, args ...interface{}) {
	e.FieldLogger.Infof(s, args...)
}

func (e *ELKLogger) Slowv(v interface{}) {
	e.FieldLogger.Info(v)
}

func (e *ELKLogger) WithDuration(d time.Duration) *ELKLogger {
	e.Duration = ReprOfDuration(d)
	return e
}

// ReprOfDuration returns the string representation of given duration in ms.
func ReprOfDuration(duration time.Duration) string {
	return fmt.Sprintf("%.1fms", float32(duration)/float32(time.Millisecond))
}
