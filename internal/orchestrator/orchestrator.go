package orchestrator

import (
	"calculator/configs"
	"calculator/internal/orchestrator/handler"
	"calculator/internal/orchestrator/scheduler"
	"calculator/internal/orchestrator/storage"
	"net/http"

	"calculator/pkg/middlewares"
)

// Orchestrator represents the orchestrator.
type Orchestrator struct {
	cfg *configs.Config
}

// NewOrchestrator creates a new instance of the Orchestrator.
func NewOrchestrator(cfg *configs.Config) (*Orchestrator, error) {

	return &Orchestrator{
		cfg: cfg,
	}, nil
}

// Router returns the router for the orchestrator.
func (o *Orchestrator) Router() http.Handler {
	mux := http.NewServeMux()
	storage := storage.NewStorage()
	sheduler := scheduler.NewScheduler(storage, o.cfg)

	handler := handler.NewHandler(sheduler, storage)
	handler.RegisterRoutes(mux)

	result := middlewares.MakeLoggingMiddleware(mux)
	result = middlewares.PanicRecoveryMiddleware(result)

	return result
}
