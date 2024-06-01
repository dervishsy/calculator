package errors

import "errors"

var (
	ErrExpressionNotFound = errors.New("expression not found")
	ErrNoTasksAvailable   = errors.New("no tasks available")
	ErrTaskNotFound       = errors.New("task not found")
	ErrExpressionExists   = errors.New("expression already exists")
)
