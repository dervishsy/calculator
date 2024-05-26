package agent

import (
	"bytes"
	"calculator/internal/shared"
	"calculator/pkg/logger"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// Worker represents a computational worker that can perform arithmetic operations.
type Worker struct {
	log             logger.Logger
	orchestratorURL string
	client          *http.Client
}

// NewWorker creates a new instance of the Worker.
func NewWorker(log logger.Logger, orchestratorURL string, client *http.Client) *Worker {
	return &Worker{
		log:             log,
		orchestratorURL: orchestratorURL,
		client:          client,
	}
}

// Run starts the worker and its main loop.
func (w *Worker) Run(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			w.doWork()
		}
	}
}

func (w *Worker) doWork() {
	task, err := w.getTask()
	if err != nil {
		w.log.Errorf("Failed to get task: %v", err)
		return
	}

	if task == nil {
		time.Sleep(1 * time.Second)
		return
	}

	result, err := w.performOperation(task)
	if err != nil {
		w.log.Errorf("Failed to perform operation: %v", err)
		return
	}

	err = w.sendResult(task.ID, result)
	if err != nil {
		w.log.Errorf("Failed to send result: %v", err)
	}
}

func (w *Worker) getTask() (*shared.Task, error) {
	url := fmt.Sprintf("%s/internal/task", w.orchestratorURL)
	resp, err := w.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		w.log.Infof("No tasks available")
		return nil, nil
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var task shared.Task
	err = json.NewDecoder(resp.Body).Decode(&task)
	if err != nil {
		return nil, err
	}
	w.log.Infof("Get task: %v", task)

	return &task, nil
}

func (w *Worker) performOperation(task *shared.Task) (float64, error) {
	var result float64

	switch task.Operation {
	case "+":
		result = task.Arg1 + task.Arg2
	case "-":
		result = task.Arg1 - task.Arg2
	case "*":
		result = task.Arg1 * task.Arg2
	case "/":
		if task.Arg2 == 0 {
			return 0, fmt.Errorf("division by zero")
		}
		result = task.Arg1 / task.Arg2
	default:
		return 0, fmt.Errorf("unknown operation: %s", task.Operation)
	}

	time.Sleep(task.OperationTime)
	return result, nil
}

func (w *Worker) sendResult(taskID string, result float64) error {
	url := fmt.Sprintf("%s/internal/task", w.orchestratorURL)
	payload := shared.TaskResult{
		ID:     taskID,
		Result: result,
	}

	w.log.Infof("Send result for task %s: %f", taskID, result)

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	resp, err := w.client.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
