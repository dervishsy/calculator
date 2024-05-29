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
	orchestratorURL string
	computingPower  int
	client          *http.Client
	workers         []*Worker
	wg              sync.WaitGroup
}

// NewAgent creates a new instance of the Agent.
func NewAgent(cfg *configs.Config) (*Agent, error) {

	agent := &Agent{
		cfg:             cfg,
		orchestratorURL: cfg.OrchestratorURL,
		computingPower:  cfg.ComputingPower,
		client:          &http.Client{},
		workers:         make([]*Worker, cfg.ComputingPower),
	}

	return agent, nil
}

// Run starts the agent and its workers.
func (a *Agent) Run(ctx context.Context) {
	logger.Info("Starting agent")

	for i := 0; i < a.computingPower; i++ {
		worker := NewWorker(a.orchestratorURL, a.client)
		a.workers[i] = worker
		a.wg.Add(1)
		go worker.Run(ctx, &a.wg)
	}

	a.wg.Wait()
	logger.Info("Agent stopped")
}
