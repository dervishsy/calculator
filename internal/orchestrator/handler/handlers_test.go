package handler

import (
	"calculator/configs"
	"calculator/internal/orchestrator/scheduler"
	"calculator/internal/orchestrator/storage"
	"calculator/pkg/logger"
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

	// Create a mock logger
	mockLogger, _ := logger.NewLogger("test")

	// Create a mock storage
	mockStorage := storage.NewStorage()

	// Create a mock scheduler
	mockScheduler := scheduler.NewScheduler(mockLogger, mockStorage, cfg)

	// Create a new handler
	handler := &Handler{
		log:       mockLogger,
		scheduler: mockScheduler,
		storage:   mockStorage,
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
	mockLogger, _ := logger.NewLogger("test")

	// Create a mock storage
	mockStorage := storage.NewStorage()

	// Create a mock scheduler
	mockScheduler := scheduler.NewScheduler(mockLogger, mockStorage, cfg)

	// Create a new handler
	handler := &Handler{
		log:       mockLogger,
		scheduler: mockScheduler,
		storage:   mockStorage,
	}

	// Call the function being tested
	handler.HandleCalculate(rr, req)

	// Check the response status code
	if rr.Code != http.StatusUnprocessableEntity {
		t.Errorf("Expected status code %d, got %d", http.StatusUnprocessableEntity, rr.Code)
	}
}
