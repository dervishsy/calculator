package entities

import (
	"time"
)

// AgentTask represents a single arithmetic operation task.
type AgentTask struct {
	ExprID        string        `json:"-"`
	ID            string        `json:"id"`
	Arg1          float64       `json:"arg1"`
	Arg2          float64       `json:"arg2"`
	Operation     string        `json:"operation"`
	OperationTime time.Duration `json:"operationTime"`
}
