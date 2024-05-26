package agent

import (
	"calculator/configs"
	"calculator/pkg/logger"
	"context"
	"net/http"
	"sync"
)

// Agent represents a computational agent that can perform arithmetic operations.
type Agent struct {
	cfg             *configs.Config
	log             logger.Logger
	orchestratorURL string
	computingPower  int
	client          *http.Client
	workers         []*Worker
	wg              sync.WaitGroup
}

// NewAgent creates a new instance of the Agent.
func NewAgent(cfg *configs.Config, log logger.Logger) (*Agent, error) {

	agent := &Agent{
		cfg:             cfg,
		log:             log,
		orchestratorURL: cfg.OrchestratorURL,
		computingPower:  cfg.ComputingPower,
		client:          &http.Client{},
		workers:         make([]*Worker, cfg.ComputingPower),
	}

	return agent, nil
}

// Run starts the agent and its workers.
func (a *Agent) Run(ctx context.Context) {
	a.log.Info("Starting agent")

	for i := 0; i < a.computingPower; i++ {
		worker := NewWorker(a.log, a.orchestratorURL, a.client)
		a.workers[i] = worker
		a.wg.Add(1)
		go worker.Run(ctx, &a.wg)
	}

	a.wg.Wait()
	a.log.Info("Agent stopped")
}
