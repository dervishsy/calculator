package handler

import (
	"calculator/internal/orchestrator/impl/memory_expression_storage"
	"calculator/internal/orchestrator/impl/memory_task_storage"
	"calculator/internal/orchestrator/use_cases/scheduler"
	"calculator/internal/shared/configs"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandleCalculate_Success(t *testing.T) {
	// Create a new request with a valid request body
	reqBody := strings.NewReader(`{"id": "1", "expression": "2+2"}`)
	req, err := http.NewRequest("POST", "/calculate", reqBody)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Create a response recorder to capture the response
	rr := httptest.NewRecorder()

	cfg := &configs.Config{
		OrchestratorURL:      "http://localhost:8080",
		ComputingPower:       2,
		TimeAdditionMS:       100,
		TimeSubtractionMS:    200,
		TimeMultiplicationMS: 300,
		TimeDivisionMS:       400,
	}

	// Create a mock storage
	mockStorage := memory_expression_storage.NewStorage()

	// Create a mock task pool
	taskPool := memory_task_storage.NewTaskPool()

	// Create a mock scheduler
	mockScheduler := scheduler.NewScheduler(mockStorage, taskPool, cfg)

	// Create a new handler
	handler := &Handler{
		scheduler: mockScheduler,
	}

	// Call the function being tested
	handler.HandleCalculate(rr, req)

	// Check the response status code
	if rr.Code != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d", http.StatusCreated, rr.Code)
	}

}

func TestHandleCalculate_DecodeError(t *testing.T) {
	// Create a new request with an invalid request body
	reqBody := strings.NewReader(`invalid`)
	req, err := http.NewRequest("POST", "/calculate", reqBody)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Create a response recorder to capture the response
	rr := httptest.NewRecorder()

	cfg := &configs.Config{
		OrchestratorURL:      "http://localhost:8080",
		ComputingPower:       2,
		TimeAdditionMS:       100,
		TimeSubtractionMS:    200,
		TimeMultiplicationMS: 300,
		TimeDivisionMS:       400,
	}

	// Create a mock logger
	// Create a mock storage
	mockStorage := memory_expression_storage.NewStorage()

	// Create a mock task pool
	taskPool := memory_task_storage.NewTaskPool()

	// Create a mock scheduler
	mockScheduler := scheduler.NewScheduler(mockStorage, taskPool, cfg)

	// Create a new handler
	handler := &Handler{
		scheduler: mockScheduler,
	}

	// Call the function being tested
	handler.HandleCalculate(rr, req)

	// Check the response status code
	if rr.Code != http.StatusUnprocessableEntity {
		t.Errorf("Expected status code %d, got %d", http.StatusUnprocessableEntity, rr.Code)
	}
}
