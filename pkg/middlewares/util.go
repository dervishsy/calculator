package middlewares

import (
	"bytes"
	"calculator/pkg/logger"
	"context"
	"io"
	"net/http"
)

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
