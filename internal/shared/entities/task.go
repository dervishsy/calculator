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

type ArgType = int

const (
	IsTask ArgType = iota
	IsNumber
)

// Task represents a task in the task pool.
type Task struct {
	ExprID    string
	ID        string
	ArgLeft   Arg
	ArgRight  Arg
	Operation string
	Result    float64
}

// Arg represents an argument in a task.
type Arg struct {
	ArgFloat float64
	ArgTask  *Task
	ArgType  ArgType
}
