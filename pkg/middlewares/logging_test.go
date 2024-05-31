package middlewares

import (
	"calculator/pkg/utils"
	"calculator/pkg/uuid"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMakeLoggingMiddleware(t *testing.T) {
	// Test case: request with no request ID in header
	t.Run("NoRequestID", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := MakeLoggingMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Context().Value(utils.HeaderRequestID).(string)
			if requestID == "" {
				t.Error("Expected request ID in context")
			}
		}))
		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, rr.Code)
		}
	})

	// Test case: request with request ID in header
	t.Run("RequestID", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/", nil)
		if err != nil {
			t.Fatal(err)
		}
		requestID := uuid.New()
		req.Header.Set(string(utils.HeaderRequestID), requestID)

		rr := httptest.NewRecorder()
		handler := MakeLoggingMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestIDl := r.Context().Value(utils.HeaderRequestID).(string)
			if requestID != requestIDl {
				t.Errorf("Expected request ID %s, got %s", requestID, r.Context().Value(utils.HeaderRequestID))
			}
		}))
		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, rr.Code)
		}
	})

	// Test case: request to private path
	// 	t.Run("PrivatePath", func(t *testing.T) {
	// 		req, err := http.NewRequest("GET", "/readyz", nil)
	// 		if err != nil {
	// 			t.Fatal(err)
	// 		}

	// 		rr := httptest.NewRecorder()
	// 		handler := MakeLoggingMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 			t.Error("Should not log request to private path")
	// 		}))
	// 		handler.ServeHTTP(rr, req)

	// 		if rr.Code != http.StatusOK {
	// 			t.Errorf("Expected status code %d, got %d", http.StatusOK, rr.Code)
	// 		}
	// 	})
}
