package middlewares

import (
	"calculator/pkg/utils"
	"calculator/pkg/uuid"
	"context"
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
