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
}
