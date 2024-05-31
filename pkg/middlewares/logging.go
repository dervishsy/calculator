package middlewares

import (
	"bytes"
	"calculator/pkg/logger"
	"calculator/pkg/utils"
	"calculator/pkg/uuid"
	"context"
	"io"
	"net/http"
)

var privatePaths map[string]struct{} = map[string]struct{}{
	"/readyz":  {},
	"/healthz": {},
	"/metrics": {},
}

func MakeLoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Header.Get(string(utils.HeaderRequestID))
		if requestID == "" {
			requestID = uuid.New()
		}

		ctx := context.WithValue(r.Context(), utils.HeaderRequestID, requestID)
		w.Header().Set(string(utils.HeaderRequestID), requestID)

		if _, exist := privatePaths[r.URL.Path]; !exist {
			defer logRequest(ctx, "finish request", r)
			logRequest(ctx, "start request", r)
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func logRequest(_ context.Context, message string, r *http.Request) {
	var body []byte
	if r.Body != nil {
		body, _ = io.ReadAll(r.Body)
	}

	r.Body = io.NopCloser(bytes.NewBuffer(body))
	logger.Info(
		message,
		"method", r.Method,
		"uri", r.URL.String(),
		"remote_addr", r.RemoteAddr,
		"user_agent", r.UserAgent(),
		"headers", r.Header,
		"request_body_size", len(body),
		"request_body", body)
}
