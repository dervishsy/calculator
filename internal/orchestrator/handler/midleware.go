package handler

import (
	"bytes"
	"io"
	"net/http"
)

func (h *Handler) LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		b, _ := io.ReadAll(r.Body)
		r.Body.Close() //  must close
		r.Body = io.NopCloser(bytes.NewBuffer(b))

		h.log.Infof("%s %s - %s", r.Method, r.URL.Path, b)
		next.ServeHTTP(w, r)
	})
}

func (h *Handler) PanicRecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				h.log.Errorf("%v", err)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
