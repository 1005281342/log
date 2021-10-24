package main

import (
	"context"
	"fmt"
	"github.com/1005281342/log"
)

func main() {
	var ctx = context.Background()
	ctx = log.SetTraceID(ctx, "jesontest20210124-3")

	var logger, err = log.WithContext(ctx,
		log.WithAddress("127.0.0.1:5000"),
		log.WithAppName("jeson"),
		log.WithFuncName("test"),
	)
	if err != nil {
		fmt.Printf("new logger failed: %+v \n", err)
		return
	}
	logger.Error("test Error")
	logger.Errorf("test Error%s", "f")
	logger.Errorv("test Errorv")
	logger.Info("test Info")
	logger.Infof("test Info%s", "f")
	logger.Infov("test Infov")
	logger.Slow("test Slow")
	logger.Slowf("test Slow%s", "f")
	logger.Slowv("test Slowv")
}
