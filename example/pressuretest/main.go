package main

import (
	"context"
	"fmt"
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
			//sendMessage()
			//sendMessageWithXELKLogger()
			//sendMessageWithXELKLogger2()
			sendMessageWithXELKLogger3()
		case <-lc.C:
			return
		}
	}
}

var (
	cnt    int32
	logger *log.ELKLogger
)

func init() {
	var ctx = context.Background()
	ctx = log.SetTraceID(ctx, "jesontest20210125-1-3")

	var err error
	logger, err = log.WithContext(ctx,
		log.WithAddress("127.0.0.1:5000"),
		log.WithAppName("jeson"),
		log.WithFuncName("pressure test"),
	)

	if err != nil {
		fmt.Println("err: %+v", err)
	}
}

func sendMessage() {
	logger.Infof("hello %d", cnt)
	atomic.AddInt32(&cnt, 1)
}

func sendMessageWithXELKLogger3() {
	var ctx = context.Background()
	ctx = log.SetTraceID(ctx, "jesontest20210126-6")
	var loggerX = log.NewXELKLoggerWithContext(ctx, log.WithAddress("127.0.0.1:5000"),
		log.WithAppName("jeson"),
		log.WithFuncName("sendMessageWithXELKLogger3"),
	)
	loggerX.Infof("hello %d", cnt)
	atomic.AddInt32(&cnt, 1)
}

//func sendMessageWithXELKLogger2() {
//	var ctx = context.Background()
//	ctx = log.SetTraceID(ctx, "jesontest20210126-5")
//	var loggerX, connCloseFunc = log.NewXELKLoggerWithContext(ctx, log.WithAddress("127.0.0.1:5000"),
//		log.WithAppName("jeson"),
//		log.WithFuncName("sendMessageWithXELKLogger"),
//	)
//	defer func() {
//		if err := connCloseFunc(); err != nil {
//			fmt.Printf("close conn failed: %+v \n", err)
//		}
//	}()
//	loggerX.Infof("hello %d", cnt)
//	atomic.AddInt32(&cnt, 1)
//}

//func sendMessageWithXELKLogger() {
//	var ctx = context.Background()
//	ctx = log.SetTraceID(ctx, "jesontest20210126-1")
//	var loggerX = log.NewXELKLoggerWithContext(ctx, log.WithAddress("127.0.0.1:5000"),
//		log.WithAppName("jeson"),
//		log.WithFuncName("sendMessageWithXELKLogger"),
//	)
//	loggerX.Infof("hello %d", cnt)
//	atomic.AddInt32(&cnt, 1)
//}
