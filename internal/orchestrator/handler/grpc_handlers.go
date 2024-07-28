package handler

import (
	"calculator/internal/orchestrator/use_cases/scheduler"
	"calculator/proto/calculator/proto"
	"context"
)

type GRPCHandler struct {
	proto.UnimplementedCalculatorServer
	scheduler *scheduler.Scheduler
}

func NewGRPCHandler(scheduler *scheduler.Scheduler) *GRPCHandler {
	return &GRPCHandler{
		scheduler: scheduler,
	}
}

func (h *GRPCHandler) GetTask(ctx context.Context, req *proto.GetTaskRequest) (*proto.Task, error) {
	task, err := h.scheduler.GetTask()
	if err != nil {
		return nil, err
	}
	return &proto.Task{
		ExprId:        task.ExprID,
		Id:            task.ID,
		Arg1:          task.Arg1,
		Arg2:          task.Arg2,
		Operation:     task.Operation,
		OperationTime: int64(task.OperationTime),
	}, nil
}

func (h *GRPCHandler) SubmitResult(ctx context.Context, result *proto.TaskResult) (*proto.SubmitResultResponse, error) {
	err := h.scheduler.ProcessResult(result.Id, result.Result)
	if err != nil {
		return nil, err
	}
	return &proto.SubmitResultResponse{}, nil
}
