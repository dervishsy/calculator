package shared

import (
	"time"
)

// ExpressionStatus represents the status of an arithmetic expression.
type ExpressionStatus string

const (
	ExpressionStatusPending    ExpressionStatus = "pending"
	ExpressionStatusProcessing ExpressionStatus = "processing"
	ExpressionStatusCompleted  ExpressionStatus = "completed"
)

// Expression represents an arithmetic expression and its current status.
type Expression struct {
	ID         string           `json:"id"`
	Expression string           `json:"expression"`
	Status     ExpressionStatus `json:"status"`
	Result     float64          `json:"result,omitempty"`
}

// Task represents a single arithmetic operation task.
type Task struct {
	ExprID        string        `json:"-"`
	ID            string        `json:"id"`
	Arg1          float64       `json:"arg1"`
	Arg2          float64       `json:"arg2"`
	Operation     string        `json:"operation"`
	OperationTime time.Duration `json:"operationTime"`
}

// TaskResult represents the result of a task computation.
type TaskResult struct {
	ID     string  `json:"id"`
	Result float64 `json:"result"`
}
