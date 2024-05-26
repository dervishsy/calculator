package handler

import (
	"calculator/internal/orchestrator/scheduler"
	"calculator/internal/orchestrator/storage"
	"calculator/internal/shared"
	"calculator/pkg/logger"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
)

// Handler represents the HTTP handler for the orchestrator.
type Handler struct {
	log       logger.Logger
	scheduler *scheduler.Scheduler
	storage   *storage.Storage
}

// NewHandler creates a new instance of the Handler.
func NewHandler(log logger.Logger, scheduler *scheduler.Scheduler, storage *storage.Storage) *Handler {
	return &Handler{
		log:       log,
		scheduler: scheduler,
		storage:   storage,
	}
}

// HandleCalculate handles the request to calculate an arithmetic expression.
func (h *Handler) HandleCalculate(w http.ResponseWriter, r *http.Request) {
	var expr shared.Expression
	err := json.NewDecoder(r.Body).Decode(&expr)
	if err != nil {
		h.log.Errorf("Failed to decode request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusUnprocessableEntity)
		return
	}
	defer r.Body.Close()

	err = h.scheduler.ScheduleExpression(expr.ID, expr.Expression)
	if err != nil {
		h.log.Errorf("Failed to schedule expression: %v", err)
		http.Error(w, fmt.Sprintf("Failed to schedule expression: %v", err), http.StatusInternalServerError)
		return
	}
	h.log.Infof("Schedule expression: %v", expr)

	w.WriteHeader(http.StatusCreated)
}

// HandleGetExpressions handles the request to get a list of expressions.
func (h *Handler) HandleGetExpressions(w http.ResponseWriter, r *http.Request) {
	expressions, err := h.storage.GetExpressions()
	if err != nil {
		h.log.Errorf("Failed to get expressions: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	respBody, err := json.Marshal(map[string][]shared.Expression{
		"expressions": expressions,
	})
	if err != nil {
		h.log.Errorf("Failed to marshal response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	h.log.Infof("Get a list of expressions: %v", expressions)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(respBody)
}

// HandleGetExpression handles the request to get a specific expression.
func (h *Handler) HandleGetExpression(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	expr, err := h.storage.GetExpression(id)
	if err != nil {
		if err == storage.ErrExpressionNotFound {
			http.Error(w, "Expression not found", http.StatusNotFound)
			return
		}
		h.log.Errorf("Failed to get expression: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	respBody, err := json.Marshal(map[string]shared.Expression{
		"expression": *expr,
	})
	if err != nil {
		h.log.Errorf("Failed to marshal response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	h.log.Infof("Request to get a specific expression: %v", expr)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(respBody)
}

// HandleGetTask handles the request to get a task for computation.
func (h *Handler) HandleGetTask(w http.ResponseWriter, r *http.Request) {
	task, err := h.scheduler.GetTask()
	if err != nil {
		if err == scheduler.ErrNoTasksAvailable {
			http.Error(w, "No tasks available", http.StatusNotFound)
			return
		}
		h.log.Errorf("Failed to get task: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	respBody, err := json.Marshal(task)
	if err != nil {
		h.log.Errorf("Failed to marshal response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	h.log.Infof("Request to get a task for computation: %v", task)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(respBody)
}

// HandlePostTask handles the request to post a result of a task computation.
func (h *Handler) HandlePostTask(w http.ResponseWriter, r *http.Request) {
	var result shared.TaskResult
	err := json.NewDecoder(r.Body).Decode(&result)
	if err != nil {
		h.log.Errorf("Failed to decode request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusUnprocessableEntity)
		return
	}
	defer r.Body.Close()

	err = h.scheduler.ProcessResult(result.ID, result.Result)
	if err != nil {
		h.log.Errorf("Failed to process result: %v", err)
		if err == scheduler.ErrTaskNotFound {
			http.Error(w, "Task not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	h.log.Infof("Request to post a result of a task computation: %v", result)

	w.WriteHeader(http.StatusOK)
}

// HandleStartPage handles the request to get the start page.
func (h *Handler) HandleStartPage(w http.ResponseWriter, r *http.Request) {
	templates := template.Must(template.ParseFiles(
		//filepath.Join("web/templates", "base.html"),
		filepath.Join("web/templates", "index.html"),
	))
	//	tmpl := "base.html"
	tmpl := "index.html"
	err := templates.ExecuteTemplate(w, tmpl, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}