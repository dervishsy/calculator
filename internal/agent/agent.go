package agent

import (
	"calculator/internal/shared/configs"
	"calculator/pkg/logger"
	"context"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Agent represents a computational agent that can perform arithmetic operations.

type Agent struct {
	cfg            *configs.Config
	computingPower int
	conn           *grpc.ClientConn
	workers        []*Worker
	wg             sync.WaitGroup
}

func NewAgent(cfg *configs.Config) (*Agent, error) {
	conn, err := grpc.NewClient(cfg.OrchestratorURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	agent := &Agent{
		cfg:            cfg,
		computingPower: cfg.ComputingPower,
		conn:           conn,
		workers:        make([]*Worker, cfg.ComputingPower),
	}

	return agent, nil
}

func (a *Agent) Run(ctx context.Context) {
	logger.Info("Starting agent")

	for i := 0; i < a.computingPower; i++ {
		worker := NewWorker(a.conn)
		a.workers[i] = worker
		a.wg.Add(1)
		go func() {
			defer a.wg.Done()
			worker.Run(ctx)
		}()
	}

	a.wg.Wait()
	logger.Info("Agent stopped")
}
