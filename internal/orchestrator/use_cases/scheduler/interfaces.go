package scheduler

import (
	"calculator/internal/shared/entities"
)

type ExpressionService interface {
	CreateExpression(id string, expr string) error
	GetExpression(id string) (*entities.Expression, error)
	GetExpressions() ([]entities.Expression, error)
	UpdateExpression(id string, status entities.ExpressionStatus, result float64) error
}

type TaskService interface {
	AddTasks(tasks []entities.Task) error
	GetTaskToCompute() (entities.Task, error)
	SetTaskResultAfterCompute(id string, result float64) error
	DeleteTask(id string) error
	DeleteExpression(id string) error
	IsLastTask(id string) (bool, error)
	GetExpressionIDByTaskID(taskID string) (string, error)
}
