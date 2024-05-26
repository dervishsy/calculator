package tasks

import (
	"fmt"
)

type ArgType = int

const (
	isTask ArgType = iota
	isNumber
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

// String returns a string representation of the task.
func (t Task) String() string {
	return fmt.Sprintf("%s: %s = %v %v %v", t.ExprID, t.ID, t.ArgLeft, t.Operation, t.ArgRight)
}

// Compare compares two tasks and returns true if they are the same.
func (t Task) Compare(t2 Task) bool {
	return t.ExprID == t2.ExprID &&
		t.Operation == t2.Operation &&
		t.ArgLeft.ArgFloat == t2.ArgLeft.ArgFloat &&
		t.ArgRight.ArgFloat == t2.ArgRight.ArgFloat
}

// Compute computes the result of the task.
func (t Task) Compute() (float64, error) {
	if t.ArgLeft.ArgType != isNumber || t.ArgRight.ArgType != isNumber {
		return 0, fmt.Errorf("invalid arguments")
	}
	switch t.Operation {
	case "+":
		return t.ArgLeft.ArgFloat + t.ArgRight.ArgFloat, nil
	case "-":
		return t.ArgLeft.ArgFloat - t.ArgRight.ArgFloat, nil
	case "*":
		return t.ArgLeft.ArgFloat * t.ArgRight.ArgFloat, nil
	case "/":
		return t.ArgLeft.ArgFloat / t.ArgRight.ArgFloat, nil
	default:
		return 0, fmt.Errorf("unknown operation: %s", t.Operation)
	}
}

// String returns a string representation of the argument.
func (a Arg) String() string {
	if a.ArgType == isNumber {
		return fmt.Sprintf("%v", a.ArgFloat)
	} else {
		return fmt.Sprintf("%v", a.ArgTask.ID)
	}
}
