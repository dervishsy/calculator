package handler

import (
	"calculator/internal/orchestrator/expression_storage"
	"calculator/internal/orchestrator/scheduler"
	"calculator/internal/shared/entities"
	"calculator/pkg/logger"
	"calculator/pkg/utils"
	"encoding/json"
	"net/http"
)

// Handler represents the HTTP handler for the orchestrator.
type Handler struct {
	scheduler *scheduler.Scheduler
	storage   *expression_storage.Storage
}

// NewHandler creates a new instance of the Handler.
func NewHandler(scheduler *scheduler.Scheduler, storage *expression_storage.Storage) *Handler {
	return &Handler{
		scheduler: scheduler,
		storage:   storage,
	}
}

// HandleCalculate handles the request to calculate an arithmetic expression.
func (h *Handler) HandleCalculate(w http.ResponseWriter, r *http.Request) {
	var expr entities.Expression
	err := json.NewDecoder(r.Body).Decode(&expr)
	if err != nil {
		logger.Errorf("Failed to decode request body: %v", err)
		if err = utils.RespondWith422(w); err != nil {
			logger.Error(err)
		}
		return
	}
	defer r.Body.Close()

	err = h.scheduler.ScheduleExpression(expr.ID, expr.Expression)
	if err != nil {
		logger.Errorf("Failed to schedule expression: %v", err)
		if err = utils.RespondWith500(w); err != nil {
			logger.Error(err)
		}
		return
	}
	logger.Infof("Schedule expression: %v", expr)

	if err = utils.SuccessRepondWith201(w, struct{}{}); err != nil {
		logger.Error(err)
	}

}

// HandleGetExpressions handles the request to get a list of expressions.
func (h *Handler) HandleGetExpressions(w http.ResponseWriter, r *http.Request) {
	expressions, err := h.storage.GetExpressions()
	if err != nil {
		logger.Errorf("Failed to get expressions: %v", err)
		if err = utils.RespondWith500(w); err != nil {
			logger.Error(err)
		}
		return
	}

	logger.Infof("Get a list of expressions: %v", expressions)

	ex := map[string][]entities.Expression{"expressions": expressions}
	if err = utils.SuccessRespondWith200(w, ex); err != nil {
		logger.Error(err)
	}

}

// HandleGetExpression handles the request to get a specific expression.
func (h *Handler) HandleGetExpression(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	expr, err := h.storage.GetExpression(id)
	if err != nil {
		if err == expression_storage.ErrExpressionNotFound {
			if err = utils.RespondWith404(w); err != nil {
				logger.Error(err)
			}
			return
		}
		logger.Errorf("Failed to get expression: %v", err)
		if err = utils.RespondWith500(w); err != nil {
			logger.Error(err)
		}
		return
	}

	logger.Infof("Request to get a specific expression: %v", expr)

	ex := map[string]entities.Expression{"expression": *expr}
	if err = utils.SuccessRespondWith200(w, ex); err != nil {
		logger.Error(err)
	}

}

// HandleGetTask handles the request to get a task for computation.
func (h *Handler) HandleGetTask(w http.ResponseWriter, r *http.Request) {
	task, err := h.scheduler.GetTask()
	if err != nil {
		if err == scheduler.ErrNoTasksAvailable {
			if err = utils.RespondWith404(w); err != nil {
				logger.Error(err)
			}
			return
		}
		logger.Errorf("Failed to get task: %v", err)
		if err = utils.RespondWith500(w); err != nil {
			logger.Error(err)
		}
		return
	}

	logger.Infof("Request to get a task for computation: %v", task)

	if err = utils.SuccessRespondWith200(w, task); err != nil {
		logger.Error(err)
	}

}

// HandlePostTask handles the request to post a result of a task computation.
func (h *Handler) HandlePostTask(w http.ResponseWriter, r *http.Request) {
	var result entities.TaskResult
	err := json.NewDecoder(r.Body).Decode(&result)
	if err != nil {
		logger.Errorf("Failed to decode request body: %v", err)
		if err = utils.RespondWith422(w); err != nil {
			logger.Error(err)
		}
		return
	}
	defer r.Body.Close()

	err = h.scheduler.ProcessResult(result.ID, result.Result)
	if err != nil {
		logger.Errorf("Failed to process result: %v", err)
		if err == scheduler.ErrTaskNotFound {
			if err = utils.RespondWith404(w); err != nil {
				logger.Error(err)
			}
			return
		}
		if err = utils.RespondWith500(w); err != nil {
			logger.Error(err)
		}
		return
	}
	logger.Infof("Request to post a result of a task computation: %v", result)

	if err = utils.SuccessRespondWith200(w, struct{}{}); err != nil {
		logger.Error(err)
	}

}
