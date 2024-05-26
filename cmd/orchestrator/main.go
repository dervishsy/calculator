package main

import (
	"calculator/configs"
	"calculator/internal/orchestrator"
	"calculator/pkg/logger"
	"fmt"
	"net/http"
	"os"
)

func main() {
	// Инициализация логгера
	log, err := logger.NewLogger("orchestrator")
	if err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}

	// Загрузка конфигурации
	cfg, err := configs.LoadConfig("configs/orchestrator.yml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Создание оркестратора
	orchestrator, err := orchestrator.NewOrchestrator(cfg, log)
	if err != nil {
		log.Fatalf("Failed to create orchestrator: %v", err)
	}

	// Запуск HTTP-сервера
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Infof("Starting server on %s", addr)
	err = http.ListenAndServe(addr, orchestrator.Router())
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
