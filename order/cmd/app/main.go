package main

import (
	"order/cmd/db"
	"order/transport/grpc"
	"sync"

	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.Info("Starting Application")

	// Initialize Redis
	if err := db.InitRedis(); err != nil {
		logrus.Warnf("Redis initialization failed: %v. Continuing without cache.", err)
	}

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		logrus.Info("Calling GRPCListen()")
		grpc.GRPCListen()
		logrus.Info("GRPCListen() exited")
	}()

	wg.Wait()
}
