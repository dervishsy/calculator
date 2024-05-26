package storage

import (
	"calculator/internal/orchestrator/tasks"
	"calculator/internal/shared"
	"errors"
	"slices"
	"sync"
)

var (
	ErrExpressionNotFound = errors.New("expression not found")
)

// Storage represents a simple in-memory storage for arithmetic expressions.
type Storage struct {
	expressions map[string]*shared.Expression
	taskIDs     map[string]string
	mu          sync.RWMutex
}

// NewStorage creates a new instance of the Storage.
func NewStorage() *Storage {
	return &Storage{
		expressions: make(map[string]*shared.Expression),
		taskIDs:     make(map[string]string),
	}
}

// CreateExpression creates a new arithmetic expression.
func (s *Storage) CreateExpression(id string, expr string, tasksList []tasks.Task) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.expressions[id]; ok {
		return errors.New("expression already exists")
	}

	s.expressions[id] = &shared.Expression{
		ID:         id,
		Expression: expr,
		Status:     shared.ExpressionStatusPending,
	}
	for _, task := range tasksList {
		s.taskIDs[task.ID] = id
	}
	return nil
}

// GetExpression retrieves an arithmetic expression by its ID.
func (s *Storage) GetExpression(id string) (*shared.Expression, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	expr, ok := s.expressions[id]
	if !ok {
		return nil, ErrExpressionNotFound
	}

	return expr, nil
}

// GetExpressions retrieves all arithmetic expressions.
func (s *Storage) GetExpressions() ([]shared.Expression, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	expressions := make([]shared.Expression, 0, len(s.expressions))
	for _, expr := range s.expressions {
		expressions = append(expressions, *expr)
	}

	slices.SortFunc(expressions, func(a, b shared.Expression) int {
		if a.ID == b.ID {
			return 0
		} else if a.ID > b.ID {
			return 1
		}
		return -1
	})

	return expressions, nil
}

// UpdateExpression updates the status and result of an arithmetic expression.
func (s *Storage) UpdateExpression(id string, status shared.ExpressionStatus, result float64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	expr, ok := s.expressions[id]
	if !ok {
		return ErrExpressionNotFound
	}

	expr.Result = result
	expr.Status = status

	return nil
}

// GetExpressionByTaskID retrieves the arithmetic expression associated with a task ID.
func (s *Storage) GetExpressionByTaskID(taskID string) (*shared.Expression, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	exprID, ok := s.taskIDs[taskID]
	if !ok {
		return nil, ErrExpressionNotFound
	}

	expr, ok := s.expressions[exprID]
	if !ok {
		return nil, ErrExpressionNotFound
	}

	return expr, nil
}
