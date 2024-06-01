package main

import (
	"calculator/internal/orchestrator"
	"calculator/internal/shared/configs"
	"calculator/pkg/logger"
	"context"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	// Configuration loading
	conf, err := configs.LoadConfig("configs/orchestrator.yml")
	if err != nil {
		logger.Fatalf("Failed to load config: %v", err)
	}

	serverCtx, serverStopCtx := context.WithCancel(context.Background())
	app, err := orchestrator.New(conf)
	if err != nil {
		logger.Error("failed to read config")
		os.Exit(1)
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		app.GracefulStop(serverCtx, sig, serverStopCtx)
	}()

	err = app.Run()
	if err != nil {
		logger.Error("failed to start app")
		os.Exit(1)
	}

	<-serverCtx.Done()

}
