package memory_expression_storage

import (
	use_cases_errors "calculator/internal/orchestrator/use_cases/errors"
	"calculator/internal/shared/entities"
	"reflect"
	"testing"
)

func TestCreateExpression(t *testing.T) {
	// Test case 1: Successfully create a new expression
	t.Run("Successfully create a new expression", func(t *testing.T) {
		storage := NewStorage()
		id := "1"
		expr := "2+2"

		err := storage.CreateExpression(id, expr)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		// Check if the expression was added to the storage
		if _, ok := storage.expressions[id]; !ok {
			t.Errorf("Expected expression with ID %s to be added to storage", id)
		}
	})

	// Test case 2: Return an error when trying to create an expression with an ID that already exists
	t.Run("Return an error when trying to create an expression with an ID that already exists", func(t *testing.T) {
		storage := NewStorage()
		id := "1"
		expr := "2+2"

		// Create the expression for the first time
		err := storage.CreateExpression(id, expr)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		// Try to create the expression again
		err = storage.CreateExpression(id, expr)
		if err == nil {
			t.Errorf("Expected error, got nil")
		}
	})
}

func TestGetExpression(t *testing.T) {
	// Test case 1: Getting an existing expression by ID.
	t.Run("Getting an existing expression by ID", func(t *testing.T) {
		storage := &Storage{
			expressions: map[string]*entities.Expression{
				"1": {
					ID:         "1",
					Expression: "2+2",
					Status:     entities.ExpressionStatusPending,
				},
			},
		}
		expr, err := storage.GetExpression("1")
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		expectedExpr := &entities.Expression{
			ID:         "1",
			Expression: "2+2",
			Status:     entities.ExpressionStatusPending,
		}
		if !reflect.DeepEqual(expr, expectedExpr) {
			t.Errorf("Expected expression %v, got %v", expectedExpr, expr)
		}
	})

	// Test case 2: Getting a non-existing expression by ID.
	t.Run("Getting a non-existing expression by ID", func(t *testing.T) {
		storage := &Storage{
			expressions: map[string]*entities.Expression{},
		}
		_, err := storage.GetExpression("1")
		if err != use_cases_errors.ErrExpressionNotFound {
			t.Errorf("Expected error %v, got %v", use_cases_errors.ErrExpressionNotFound, err)
		}
	})
}

func TestGetExpressions(t *testing.T) {
	// Test case 1: Getting all expressions when there are multiple expressions.
	t.Run("Getting all expressions when there are multiple expressions", func(t *testing.T) {
		storage := &Storage{
			expressions: map[string]*entities.Expression{
				"1": {
					ID:         "1",
					Expression: "2+2",
					Status:     entities.ExpressionStatusPending,
				},
				"2": {
					ID:         "2",
					Expression: "3+3",
					Status:     entities.ExpressionStatusPending,
				},
			},
		}
		expressions, err := storage.GetExpressions()
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		expectedExpressions := []entities.Expression{
			{
				ID:         "1",
				Expression: "2+2",
				Status:     entities.ExpressionStatusPending,
			},
			{
				ID:         "2",
				Expression: "3+3",
				Status:     entities.ExpressionStatusPending,
			},
		}
		if !reflect.DeepEqual(expressions, expectedExpressions) {
			t.Errorf("Expected expressions %v, got %v", expectedExpressions, expressions)
		}
	})

	// Test case 2: Getting an empty slice when there are no expressions.
	t.Run("Getting an empty slice when there are no expressions", func(t *testing.T) {
		storage := &Storage{
			expressions: map[string]*entities.Expression{},
		}
		expressions, err := storage.GetExpressions()
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if len(expressions) != 0 {
			t.Errorf("Expected empty slice, got %v", expressions)
		}
	})

	// Test case 3: Getting expressions sorted by ID in ascending order.
	t.Run("Getting expressions sorted by ID in ascending order", func(t *testing.T) {
		storage := &Storage{
			expressions: map[string]*entities.Expression{
				"2": {
					ID:         "2",
					Expression: "3+3",
					Status:     entities.ExpressionStatusPending,
				},
				"1": {
					ID:         "1",
					Expression: "2+2",
					Status:     entities.ExpressionStatusPending,
				},
			},
		}
		expressions, err := storage.GetExpressions()
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		expectedExpressions := []entities.Expression{
			{
				ID:         "1",
				Expression: "2+2",
				Status:     entities.ExpressionStatusPending,
			},
			{
				ID:         "2",
				Expression: "3+3",
				Status:     entities.ExpressionStatusPending,
			},
		}
		if !reflect.DeepEqual(expressions, expectedExpressions) {
			t.Errorf("Expected expressions %v, got %v", expectedExpressions, expressions)
		}

	})
}

func TestUpdateExpression(t *testing.T) {
	// Test case: expression exists and status and result are updated
	storage := &Storage{
		expressions: map[string]*entities.Expression{
			"1": {
				ID:         "1",
				Expression: "1+1",
				Status:     entities.ExpressionStatusPending,
				Result:     0,
			},
		},
	}
	err := storage.UpdateExpression("1", entities.ExpressionStatusProcessing, 2)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if storage.expressions["1"].Status != entities.ExpressionStatusProcessing {
		t.Errorf("expected status to be %v, got %v", entities.ExpressionStatusProcessing, storage.expressions["1"].Status)
	}
	if storage.expressions["1"].Result != 2 {
		t.Errorf("expected result to be %v, got %v", 2, storage.expressions["1"].Result)
	}

	// Test case: expression does not exist
	err = storage.UpdateExpression("2", entities.ExpressionStatusProcessing, 2)
	if err != use_cases_errors.ErrExpressionNotFound {
		t.Errorf("expected error to be %v, got %v", use_cases_errors.ErrExpressionNotFound, err)
	}
}
