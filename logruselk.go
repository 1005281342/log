package log

import (
	"context"
	"fmt"
	"math/rand"
	"net"
	"sync"
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
		return logger
	}
	logger = &XELKLogger{ELKLogger: loggerELK}
	return logger
}

//// NewXELKLoggerWithContext new go-zero elk logger
//func NewXELKLoggerWithContext(ctx context.Context, optFuncList ...OptionFunc) (logx.Logger, func() error) {
//	var (
//		logger         logx.Logger
//		loggerELK, err = WithContext(ctx, optFuncList...)
//	)
//	if err != nil {
//		logger = logx.WithContext(ctx)
//		logger.Errorf("new ELKLogger WithContext err: %+v", err)
//		return logger, func() error {
//			return nil
//		}
//	}
//	logger = &XELKLogger{ELKLogger: loggerELK}
//	return logger, loggerELK.Close
//}

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
	//conn     net.Conn
	*Option
}

var addrConnObj = &addrConn{connPoolMap: new(sync.Map)}

type addrConn struct {
	connPoolMap *sync.Map
}

var defaultPoolSize = 32

func (a *addrConn) GetConn(addr string) (net.Conn, error) {
	var v, ok = a.connPoolMap.Load(addr)
	if !ok {
		var conn, err = net.DialTimeout("tcp", addr, 100*time.Millisecond)
		if err != nil {
			return nil, err
		}
		var connPool = make([]net.Conn, 0, defaultPoolSize)
		connPool = append(connPool, conn)
		a.connPoolMap.Store(addr, connPool)
		return conn, nil
	}

	var connPool []net.Conn
	if connPool, ok = v.([]net.Conn); !ok {
		return nil, fmt.Errorf("connPool not is `[]net.Conn`")
	}
	if len(connPool) <= defaultPoolSize {
		var conn, err = net.DialTimeout("tcp", addr, 100*time.Millisecond)
		if err != nil {
			return nil, err
		}
		connPool = append(connPool, conn)
		a.connPoolMap.Store(addr, connPool)
		return conn, nil
	}
	return connPool[rand.Intn(len(connPool))], nil
}

// GetConn 获取连接
//func (e *ELKLogger) GetConn() net.Conn {
//	return e.conn
//}

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

	if conn, err = addrConnObj.GetConn(elk.Address); err != nil {
		return nil, err
	}

	//if conn, err = net.DialTimeout("tcp", elk.Address, 100*time.Millisecond); err != nil {
	//	return nil, err
	//}
	//elk.conn = conn

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

// Close 关闭连接
//func (e *ELKLogger) Close() error {
//	if e == nil {
//		return nil
//	}
//
//	if e.conn == nil {
//		return nil
//	}
//	return e.conn.Close()
//}

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
