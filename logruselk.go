package log

import (
	"context"
	"fmt"
	"math/rand"
	"net"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/bshuster-repo/logrus-logstash-hook"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/tal-tech/go-zero/core/logx"
)

const (
	// TraceIDKey traceID key
	TraceIDKey = "traceID"
)

// NewGoZeroELKLoggerWithContext new go-zero elk logger
func NewGoZeroELKLoggerWithContext(ctx context.Context, optFuncList ...OptionFunc) logx.Logger {
	var (
		logger         logx.Logger
		loggerELK, err = WithContext(ctx, optFuncList...)
	)
	if err != nil {
		logger = logx.WithContext(ctx)
		logger.Errorf("new ELKLogger WithContext err: %+v", err)
		return logger
	}
	logger = &GoZeroELKLogger{ELKLogger: loggerELK}
	return logger
}

// GoZeroELKLogger go-zero elk logger
type GoZeroELKLogger struct {
	*ELKLogger
}

// WithDuration 打印日志携带耗时
func (g *GoZeroELKLogger) WithDuration(d time.Duration) logx.Logger {
	g.ELKLogger.WithDuration(d)
	return g
}

// ELKLogger elk logger
type ELKLogger struct {
	logrus.FieldLogger
	Duration string
	*Option
}

var addrConnObj = &addrConn{connPoolMap: new(sync.Map)}

type addrConn struct {
	connPoolMap *sync.Map
}

var defaultPoolSize = 32

func init() {
	const loadKey = "LogPoolSize"
	var s = os.Getenv(loadKey)
	if len(s) == 0 {
		return
	}
	var i, err = strconv.Atoi(s)
	if err != nil {
		return
	}
	if i > 0 {
		defaultPoolSize = i
	}
	fmt.Println(defaultPoolSize)
}

// getConn 获取TCP连接
func (a *addrConn) getConn(addr string) (net.Conn, error) {
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

// WithContext new elk logger with context
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

	if conn, err = addrConnObj.getConn(elk.Address); err != nil {
		return nil, err
	}

	var hook = logrustash.New(conn, logrustash.DefaultFormatter(fields))
	logger.Hooks.Add(hook)

	elk.fixTraceID(ctx)
	var entry = logger.WithFields(logrus.Fields{
		TraceIDKey: elk.TraceID,
	})
	elk.FieldLogger = entry
	return elk, nil
}

// WithContextAndAddress new elk logger with context and address
func WithContextAndAddress(ctx context.Context, addr string) (*ELKLogger, error) {
	return WithContext(ctx, WithAddress(addr))
}

func (e *ELKLogger) hasTraceID() bool {
	return e.TraceID != ""
}

func (e *ELKLogger) fixTraceID(ctx context.Context) {
	if e.hasTraceID() {
		return
	}
	e.TraceID = traceIdFromContext(ctx)
	if e.hasTraceID() {
		return
	}
	e.TraceID = uuid.NewString()
}

// Errorv error v
func (e *ELKLogger) Errorv(v interface{}) {
	e.FieldLogger.Error(v)
}

// Infov info v
func (e *ELKLogger) Infov(v interface{}) {
	e.FieldLogger.Info(v)
}

// Slow 目前还是通过info打印
func (e *ELKLogger) Slow(args ...interface{}) {
	e.FieldLogger.Info(args...)
}

// Slowf 目前还是通过infof打印
func (e *ELKLogger) Slowf(s string, args ...interface{}) {
	e.FieldLogger.Infof(s, args...)
}

// Slowv 目前还是通过infov打印
func (e *ELKLogger) Slowv(v interface{}) {
	e.FieldLogger.Info(v)
}

// WithDuration 打印日志携带耗时
func (e *ELKLogger) WithDuration(d time.Duration) *ELKLogger {
	e.Duration = ReprOfDuration(d)
	return e
}

// ReprOfDuration returns the string representation of given duration in ms.
func ReprOfDuration(duration time.Duration) string {
	return fmt.Sprintf("%.1fms", float32(duration)/float32(time.Millisecond))
}
