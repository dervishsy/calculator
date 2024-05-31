package storage

import (
	"calculator/internal/orchestrator/tasks"
	"calculator/internal/shared/entities"
	"errors"
	"slices"
	"sync"
)

var (
	ErrExpressionNotFound = errors.New("expression not found")
)

// Storage represents a simple in-memory storage for arithmetic expressions.
type Storage struct {
	expressions map[string]*entities.Expression
	taskIDs     map[string]string
	mu          sync.RWMutex
}

// NewStorage creates a new instance of the Storage.
func NewStorage() *Storage {
	return &Storage{
		expressions: make(map[string]*entities.Expression),
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

	s.expressions[id] = &entities.Expression{
		ID:         id,
		Expression: expr,
		Status:     entities.ExpressionStatusPending,
	}
	for _, task := range tasksList {
		s.taskIDs[task.ID] = id
	}
	return nil
}

// GetExpression retrieves an arithmetic expression by its ID.
func (s *Storage) GetExpression(id string) (*entities.Expression, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	expr, ok := s.expressions[id]
	if !ok {
		return nil, ErrExpressionNotFound
	}

	return expr, nil
}

// GetExpressions retrieves all arithmetic expressions.
func (s *Storage) GetExpressions() ([]entities.Expression, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	expressions := make([]entities.Expression, 0, len(s.expressions))
	for _, expr := range s.expressions {
		expressions = append(expressions, *expr)
	}

	slices.SortFunc(expressions, func(a, b entities.Expression) int {
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
func (s *Storage) UpdateExpression(id string, status entities.ExpressionStatus, result float64) error {
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
func (s *Storage) GetExpressionByTaskID(taskID string) (*entities.Expression, error) {
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
