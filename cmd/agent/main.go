package main

import (
	"calculator/internal/agent"
	"calculator/internal/shared/configs"
	"calculator/pkg/logger"
	"context"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	// Configuration loading
	cfg, err := configs.LoadConfig("configs/agent.yml")
	if err != nil {
		logger.Fatalf("Failed to load config: %v", err)
	}

	// Agent initialization
	agent, err := agent.NewAgent(cfg)
	if err != nil {
		logger.Fatalf("Failed to create agent: %v", err)
	}

	// Start agent
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go agent.Run(ctx)

	logger.Infof("Server adress must be: %s", cfg.OrchestratorURL)

	// Graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
}
