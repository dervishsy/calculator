package handler

import (
	"net/http"
)

// RegisterRoutes registers the HTTP routes for the orchestrator.
func (h *Handler) RegisterRoutes(r *http.ServeMux) {
	//api
	r.HandleFunc("POST /api/v1/calculate", h.HandleCalculate)
	r.HandleFunc("GET /api/v1/expressions/", h.HandleGetExpressions)
	r.HandleFunc("GET /api/v1/expressions/{id}/", h.HandleGetExpression)
	// agent
	r.HandleFunc("GET /internal/task", h.HandleGetTask)
	r.HandleFunc("POST /internal/task", h.HandlePostTask)
	// web
	r.HandleFunc("GET /{$}", h.HandleStartPage)
	r.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))

}
