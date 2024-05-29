package main

import (
	"calculator/configs"
	"calculator/internal/orchestrator"
	"calculator/pkg/logger"
	"fmt"
	"net/http"
)

func main() {

	// Configuration loading
	cfg, err := configs.LoadConfig("configs/orchestrator.yml")
	if err != nil {
		logger.Fatalf("Failed to load config: %v", err)
	}

	// Orchestrator initialization
	orchestrator, err := orchestrator.NewOrchestrator(cfg)
	if err != nil {
		logger.Fatalf("Failed to create orchestrator: %v", err)
	}

	// Start HTTP server
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	logger.Infof("Starting server on %s", addr)
	err = http.ListenAndServe(addr, orchestrator.Router())
	if err != nil {
		logger.Fatalf("Failed to start server: %v", err)
	}
}
