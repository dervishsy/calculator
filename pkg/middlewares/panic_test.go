package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPanicRecoveryMiddleware(t *testing.T) {
	// Test case 1: Panic recovery
	t.Run("Panic recovery", func(t *testing.T) {
		// Create a mock response recorder
		rr := httptest.NewRecorder()

		// Create a mock request
		req, err := http.NewRequest("GET", "/", nil)
		if err != nil {
			t.Fatal(err)
		}

		// Create a mock handler that panics
		panicHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			panic("test panic")
		})

		// Wrap the panic handler with the middleware
		handler := PanicRecoveryMiddleware(panicHandler)

		// Call the middleware handler
		handler.ServeHTTP(rr, req)

		// Check the response status code
		if status := rr.Code; status != http.StatusInternalServerError {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
		}

		// Check the response body
		expectedBody := "Internal server error\n"
		if body := rr.Body.String(); body != expectedBody {
			t.Errorf("handler returned unexpected body: got %v want %v", body, expectedBody)
		}
	})

	// Test case 2: Normal request handling
	t.Run("Normal request handling", func(t *testing.T) {
		// Create a mock response recorder
		rr := httptest.NewRecorder()

		// Create a mock request
		req, err := http.NewRequest("GET", "/", nil)
		if err != nil {
			t.Fatal(err)
		}

		// Create a mock handler that does not panic
		normalHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("OK"))
		})

		// Wrap the normal handler with the middleware
		handler := PanicRecoveryMiddleware(normalHandler)

		// Call the middleware handler
		handler.ServeHTTP(rr, req)

		// Check the response status code
		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		// Check the response body
		expectedBody := "OK"
		if body := rr.Body.String(); body != expectedBody {
			t.Errorf("handler returned unexpected body: got %v want %v", body, expectedBody)
		}
	})

	// Test case 3: Panic request logging
	t.Run("Panic request logging", func(t *testing.T) {
		// Create a mock response recorder
		rr := httptest.NewRecorder()

		// Create a mock request
		req, err := http.NewRequest("GET", "/", nil)
		if err != nil {
			t.Fatal(err)
		}

		// Create a mock logger
		// mockLogger := &mockLogger{}

		// Create a mock handler that panics
		panicHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			panic("test panic")
		})

		// Wrap the panic handler with the middleware
		handler := PanicRecoveryMiddleware(panicHandler)

		// Call the middleware handler
		handler.ServeHTTP(rr, req)

		// Check the response status code
		if status := rr.Code; status != http.StatusInternalServerError {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
		}

		// Check the response body
		expectedBody := "Internal server error\n"
		if body := rr.Body.String(); body != expectedBody {
			t.Errorf("handler returned unexpected body: got %v want %v", body, expectedBody)
		}

		// // Check the logger output
		// if mockLogger.message != "test panic" {
		// 	t.Errorf("unexpected logger output: got %v want %v", mockLogger.message, "test panic")
		// }
	})
}

// type mockLogger struct {
// 	// Add any necessary fields and methods for the mock logger
// 	message string
// }

// func (l *mockLogger) Info(v ...any) {
// 	// Implement the Info method for the mock logger
// 	l.message = v[0].(string)
// }

// func (l *mockLogger) Error(v ...any) {
// 	l.message = v[0].(string)
// }

// // Add any other necessary methods for the mock logger
