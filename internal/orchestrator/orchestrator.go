package orchestrator

import (
	"calculator/configs"
	"calculator/internal/orchestrator/handler"
	"calculator/internal/orchestrator/scheduler"
	"calculator/internal/orchestrator/storage"
	"calculator/pkg/logger"
	"net/http"
)

type Orchestrator struct {
	cfg *configs.Config
	log logger.Logger
}

func NewOrchestrator(cfg *configs.Config, log logger.Logger) (*Orchestrator, error) {

	return &Orchestrator{
		cfg: cfg,
		log: log,
	}, nil
}

// func (o *Orchestrator) Router() *http.ServeMux {
func (o *Orchestrator) Router() http.Handler {
	mux := http.NewServeMux()
	storage := storage.NewStorage()
	sheduler := scheduler.NewScheduler(o.log, storage, o.cfg)

	handler := handler.NewHandler(o.log, sheduler, storage)
	handler.RegisterRoutes(mux)

	result := handler.LoggingMiddleware(mux)
	result = handler.PanicRecoveryMiddleware(result)

	return result
}
