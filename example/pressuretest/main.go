package main

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/1005281342/log"
)

func main() {
	pressure(time.Second*60, 1*time.Millisecond)
}

func pressure(lifeCycle time.Duration, interval time.Duration) {
	var (
		it = time.NewTicker(interval)
		lc = time.NewTicker(lifeCycle)
	)
	defer func() {
		it.Stop()
		lc.Stop()
	}()

	for {
		select {
		case <-it.C:
			sendMessageWithGoZeroELKLogger3()
		case <-lc.C:
			return
		}
	}
}

var cnt int32

func sendMessageWithGoZeroELKLogger3() {
	var ctx = context.Background()
	var loggerX = log.NewGoZeroELKLoggerWithContext(ctx, log.WithAddress("127.0.0.1:5000"),
		log.WithAppName("jeson"),
		log.WithFuncName("sendMessageWithGoZeroELKLogger3"),
		log.WithTraceID("jesontest20211031-2"),
	)
	loggerX.Infof("hello %d", cnt)
	atomic.AddInt32(&cnt, 1)
}
