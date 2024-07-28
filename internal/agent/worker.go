package agent

import (
	"calculator/pkg/logger"
	"calculator/proto/calculator/proto"
	"context"
	"fmt"
	"strings"
	"time"

	"google.golang.org/grpc"
)

// Worker represents a computational worker that can perform arithmetic operations.
type Worker struct {
	client proto.CalculatorClient
}

// NewWorker creates a new instance of the Worker.
func NewWorker(conn *grpc.ClientConn) *Worker {
	return &Worker{
		client: proto.NewCalculatorClient(conn),
	}
}

// Run starts the worker and its main loop.
func (w *Worker) Run(ctx context.Context) {

	for {
		select {
		case <-ctx.Done():
			return
		default:
			w.doWork(ctx)
		}
	}
}

func (w *Worker) doWork(ctx context.Context) {
	task, err := w.getTask(ctx)
	if err != nil {
		if strings.Contains(err.Error(), "no tasks available") {
			time.Sleep(1 * time.Second)
			return
		}
		logger.Errorf("Failed to get task: %v", err)
		return

	}

	result, err := w.performOperation(task)
	if err != nil {
		logger.Errorf("Failed to perform operation: %v", err)
		return
	}

	err = w.sendResult(ctx, task.Id, result)
	if err != nil {
		logger.Errorf("Failed to send result: %v", err)
	}
}

func (w *Worker) getTask(ctx context.Context) (*proto.Task, error) {
	task, err := w.client.GetTask(ctx, &proto.GetTaskRequest{})
	if err != nil {
		return nil, err
	}
	logger.Infof("Get task: %v", task)
	return task, nil
}

func (w *Worker) performOperation(task *proto.Task) (float64, error) {
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

	time.Sleep(time.Duration(task.OperationTime))
	return result, nil
}

func (w *Worker) sendResult(ctx context.Context, taskID string, result float64) error {
	_, err := w.client.SubmitResult(ctx, &proto.TaskResult{
		Id:     taskID,
		Result: result,
	})
	if err != nil {
		return err
	}
	logger.Infof("Send result for task %s: %f", taskID, result)
	return nil
}
