package middlewares

import (
	"net/http"
)

func PanicRecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				logRequest(r.Context(), "Panic recovery", r)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
