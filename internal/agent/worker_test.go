package agent

import (
	"bytes"
	"calculator/internal/shared/entities"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"testing"
	"time"
)

func TestGetTask(t *testing.T) {
	// Test successful request
	worker := &Worker{
		orchestratorURL: "http://example.com",
		client:          &http.Client{},
	}

	testCases := []struct {
		name     string
		response *http.Response
		expected *entities.AgentTask
		errMsg   string
	}{
		{
			name: "valid task",
			response: &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString(`{"id": "task1", "arg1": 1, "arg2": 2, "operation": "+", "operationTime": 1}`)),
			},
			expected: &entities.AgentTask{
				ID:            "task1",
				Arg1:          1,
				Arg2:          2,
				Operation:     "+",
				OperationTime: time.Duration(1),
			},
		},
		{
			name: "no tasks available",
			response: &http.Response{
				StatusCode: http.StatusNotFound,
				Body:       io.NopCloser(bytes.NewBufferString("")),
			},
			expected: nil,
		},
		{
			name: "unexpected status code",
			response: &http.Response{
				StatusCode: http.StatusInternalServerError,
				Body:       io.NopCloser(bytes.NewBufferString("")),
			},
			errMsg: "Unexpected status code: 500",
		},
		{
			name: "JSON decoding error",
			response: &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString(`{"id": "task1", "arg1": "not a number", "arg2": 2, "operation": "+", "operationTime": "1s"}`)),
			},
			errMsg: "Failed to decode response",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			worker.client.Transport = &mockGetTransport{resp: tc.response}
			task, err := worker.getTask()

			if tc.errMsg != "" {
				if err == nil {
					t.Error("Expected error, got nil")
				}
				if err.Error() != tc.errMsg {
					t.Errorf("Expected error %q, got %q", tc.errMsg, err.Error())
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if !reflect.DeepEqual(task, tc.expected) {
				t.Errorf("Expected task %v, got %v", tc.expected, task)
			}
		})
	}
}

func TestPerformOperation(t *testing.T) {
	worker := &Worker{
		orchestratorURL: "http://example.com",
		client:          &http.Client{},
	}

	testCases := []struct {
		name     string
		task     *entities.AgentTask
		expected float64
	}{
		{
			name: "addition",
			task: &entities.AgentTask{
				Operation: "+",
				Arg1:      2,
				Arg2:      3,
			},
			expected: 5,
		},
		{
			name: "subtraction",
			task: &entities.AgentTask{
				Operation: "-",
				Arg1:      5,
				Arg2:      3,
			},
			expected: 2,
		},
		{
			name: "multiplication",
			task: &entities.AgentTask{
				Operation: "*",
				Arg1:      2,
				Arg2:      4,
			},
			expected: 8,
		},
		{
			name: "division",
			task: &entities.AgentTask{
				Operation: "/",
				Arg1:      6,
				Arg2:      3,
			},
			expected: 2,
		},
		{
			name: "division by zero",
			task: &entities.AgentTask{
				Operation: "/",
				Arg1:      6,
				Arg2:      0,
			},
		},
		{
			name: "unknown operation",
			task: &entities.AgentTask{
				Operation: "^",
				Arg1:      2,
				Arg2:      3,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := worker.performOperation(tc.task)
			if err != nil && tc.expected != 0 {
				t.Errorf("Unexpected error: %v", err)
			}
			if result != tc.expected {
				t.Errorf("Expected %f, got %f", tc.expected, result)
			}
		})
	}
}

func TestSendResult(t *testing.T) {
	tests := []struct {
		name            string
		taskID          string
		result          float64
		orchestratorURL string
		expectedErr     error
	}{
		{
			name:            "Successful send",
			taskID:          "task123",
			result:          42.0,
			orchestratorURL: "http://localhost:8080",
			expectedErr:     nil,
		},
		{
			name:            "Post error",
			taskID:          "task123",
			result:          42.0,
			orchestratorURL: "http://localhost:8080",
			expectedErr:     fmt.Errorf("Post \"http://localhost:8080/internal/task\": dial tcp: lookup localhost: No such host"),
		},
		{
			name:            "Marshal error",
			taskID:          "task123",
			result:          42.0,
			orchestratorURL: "http://localhost:8080",
			expectedErr:     fmt.Errorf(`Post "http://localhost:8080/internal/task": Failed to marshal`),
		},
		{
			name:            "Unexpected status code",
			taskID:          "task123",
			result:          42.0,
			orchestratorURL: "http://localhost:8080",
			expectedErr:     fmt.Errorf("Post \"http://localhost:8080/internal/task\": unexpected status code: 404"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &Worker{
				orchestratorURL: tt.orchestratorURL,
				client: &http.Client{
					Transport: &mockPostTransport{
						statusCode: http.StatusOK,
						err:        tt.expectedErr,
					},
				},
			}

			err := w.sendResult(tt.taskID, tt.result)

			if tt.expectedErr != nil {
				if err == nil {
					t.Errorf("Expected error: %v, got: nil", tt.expectedErr)
				} else if err.Error() != tt.expectedErr.Error() {
					t.Logf("Expected error: %v, got: %v", tt.expectedErr, err)
				}
			} else if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

type mockGetTransport struct {
	resp *http.Response
}

func (t *mockGetTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return t.resp, nil
}

type mockPostTransport struct {
	statusCode int
	err        error
}

func (t *mockPostTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return &http.Response{
		StatusCode: t.statusCode,
	}, t.err
}
