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
	AddTasks(tasks []entities.Task)
	GetTaskToCompute() (entities.Task, error)
	SetTaskResultAfterCompute(id string, result float64) error
	DeleteTask(id string)
	DeleteExpression(id string)
	IsLastTask(id string) (bool, error)
	GetExpressionIDByTaskID(taskID string) (string, error)
}
