package agent

import (
	"calculator/proto/calculator/proto"
	"context"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type mockCalculatorClient struct {
	getTaskFunc      func(ctx context.Context, in *proto.GetTaskRequest, opts ...grpc.CallOption) (*proto.Task, error)
	submitResultFunc func(ctx context.Context, in *proto.TaskResult, opts ...grpc.CallOption) (*proto.SubmitResultResponse, error)
}

func (m *mockCalculatorClient) GetTask(ctx context.Context, in *proto.GetTaskRequest, opts ...grpc.CallOption) (*proto.Task, error) {
	return m.getTaskFunc(ctx, in, opts...)
}

func (m *mockCalculatorClient) SubmitResult(ctx context.Context, in *proto.TaskResult, opts ...grpc.CallOption) (*proto.SubmitResultResponse, error) {
	return m.submitResultFunc(ctx, in, opts...)
}

func TestGetTask(t *testing.T) {
	testCases := []struct {
		name     string
		mockResp *proto.Task
		mockErr  error
		expected *proto.Task
		errMsg   string
	}{
		{
			name: "valid task",
			mockResp: &proto.Task{
				Id:            "task1",
				Arg1:          1,
				Arg2:          2,
				Operation:     "+",
				OperationTime: 1000000, // 1ms in nanoseconds
			},
			expected: &proto.Task{
				Id:            "task1",
				Arg1:          1,
				Arg2:          2,
				Operation:     "+",
				OperationTime: 1000000,
			},
		},
		{
			name:    "no tasks available",
			mockErr: status.Error(codes.NotFound, "No tasks available"),
			errMsg:  "rpc error: code = NotFound desc = No tasks available",
		},
		{
			name:    "unexpected error",
			mockErr: status.Error(codes.Internal, "Internal error"),
			errMsg:  "rpc error: code = Internal desc = Internal error",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			worker := &Worker{
				client: &mockCalculatorClient{
					getTaskFunc: func(ctx context.Context, in *proto.GetTaskRequest, opts ...grpc.CallOption) (*proto.Task, error) {
						return tc.mockResp, tc.mockErr
					},
				},
			}

			ctx := context.Background()
			task, err := worker.getTask(ctx)

			if tc.errMsg != "" {
				if err == nil {
					t.Error("Expected error, got nil")
				} else if err.Error() != tc.errMsg {
					t.Errorf("Expected error %q, got %q", tc.errMsg, err.Error())
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if task.Id != tc.expected.Id ||
				task.Arg1 != tc.expected.Arg1 ||
				task.Arg2 != tc.expected.Arg2 ||
				task.Operation != tc.expected.Operation ||
				task.OperationTime != tc.expected.OperationTime {
				t.Errorf("Expected task %v, got %v", tc.expected, task)
			}
		})
	}
}

func TestPerformOperation(t *testing.T) {
	testCases := []struct {
		name     string
		task     *proto.Task
		expected float64
		errMsg   string
	}{
		{
			name: "addition",
			task: &proto.Task{
				Operation: "+",
				Arg1:      2,
				Arg2:      3,
			},
			expected: 5,
		},
		{
			name: "subtraction",
			task: &proto.Task{
				Operation: "-",
				Arg1:      5,
				Arg2:      3,
			},
			expected: 2,
		},
		{
			name: "multiplication",
			task: &proto.Task{
				Operation: "*",
				Arg1:      2,
				Arg2:      4,
			},
			expected: 8,
		},
		{
			name: "division",
			task: &proto.Task{
				Operation: "/",
				Arg1:      6,
				Arg2:      3,
			},
			expected: 2,
		},
		{
			name: "division by zero",
			task: &proto.Task{
				Operation: "/",
				Arg1:      6,
				Arg2:      0,
			},
			errMsg: "division by zero",
		},
		{
			name: "unknown operation",
			task: &proto.Task{
				Operation: "^",
				Arg1:      2,
				Arg2:      3,
			},
			errMsg: "unknown operation: ^",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			worker := &Worker{}
			result, err := worker.performOperation(tc.task)
			if tc.errMsg != "" {
				if err == nil {
					t.Error("Expected error, got nil")
				} else if err.Error() != tc.errMsg {
					t.Errorf("Expected error %q, got %q", tc.errMsg, err.Error())
				}
				return
			}
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if result != tc.expected {
				t.Errorf("Expected %f, got %f", tc.expected, result)
			}
		})
	}
}

func TestSendResult(t *testing.T) {
	testCases := []struct {
		name        string
		taskID      string
		result      float64
		mockErr     error
		expectedErr string
	}{
		{
			name:   "Successful send",
			taskID: "task123",
			result: 42.0,
		},
		{
			name:        "Submit error",
			taskID:      "task123",
			result:      42.0,
			mockErr:     status.Error(codes.Internal, "Internal error"),
			expectedErr: "rpc error: code = Internal desc = Internal error",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			worker := &Worker{
				client: &mockCalculatorClient{
					submitResultFunc: func(ctx context.Context, in *proto.TaskResult, opts ...grpc.CallOption) (*proto.SubmitResultResponse, error) {
						return &proto.SubmitResultResponse{}, tc.mockErr
					},
				},
			}

			ctx := context.Background()
			err := worker.sendResult(ctx, tc.taskID, tc.result)

			if tc.expectedErr != "" {
				if err == nil {
					t.Error("Expected error, got nil")
				} else if err.Error() != tc.expectedErr {
					t.Errorf("Expected error %q, got %q", tc.expectedErr, err.Error())
				}
			} else if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}
