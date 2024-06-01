package memory_expression_storage

import (
	use_cases_errors "calculator/internal/orchestrator/use_cases/errors"
	"calculator/internal/shared/entities"
	"slices"
	"sync"
)

// Storage represents a simple in-memory storage for arithmetic expressions.
type Storage struct {
	expressions map[string]*entities.Expression
	mu          sync.RWMutex
}

// NewStorage creates a new instance of the Storage.
func NewStorage() *Storage {
	return &Storage{
		expressions: make(map[string]*entities.Expression),
	}
}

// CreateExpression creates a new arithmetic expression.
func (s *Storage) CreateExpression(id string, expr string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.expressions[id]; ok {
		return use_cases_errors.ErrExpressionExists
	}

	s.expressions[id] = &entities.Expression{
		ID:         id,
		Expression: expr,
		Status:     entities.ExpressionStatusPending,
	}
	return nil
}

// GetExpression retrieves an arithmetic expression by its ID.
func (s *Storage) GetExpression(id string) (*entities.Expression, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	expr, ok := s.expressions[id]
	if !ok {
		return nil, use_cases_errors.ErrExpressionNotFound
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
		return use_cases_errors.ErrExpressionNotFound
	}

	expr.Result = result
	expr.Status = status

	return nil
}
