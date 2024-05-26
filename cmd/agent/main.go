package main

import (
	"calculator/configs"
	"calculator/internal/agent"
	"calculator/pkg/logger"
	"context"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Инициализация логгера
	log, err := logger.NewLogger("agent")
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}

	// Загрузка конфигурации
	cfg, err := configs.LoadConfig("configs/agent.yml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Создание агента
	agent, err := agent.NewAgent(cfg, log)
	if err != nil {
		log.Fatalf("Failed to create agent: %v", err)
	}

	// Запуск агента
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go agent.Run(ctx)

	log.Infof("Server adress must be: %s", cfg.OrchestratorURL)

	// Обработка сигналов для корректного завершения
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
}
