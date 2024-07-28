package sqlite_expression_storage

import (
	"calculator/internal/orchestrator/impl/sqlite"
	"calculator/internal/shared/entities"
	"database/sql"
	"fmt"
)

type Storage struct {
	db *sqlite.SQLiteDB
}

func NewStorage(db *sqlite.SQLiteDB) *Storage {
	return &Storage{db: db}
}

func (s *Storage) CreateExpression(id string, expr string) error {
	_, err := s.db.Exec("INSERT INTO expressions (id, expression, status, result) VALUES (?, ?, ?, ?)",
		id, expr, entities.ExpressionStatusPending, 0)
	return err
}

func (s *Storage) GetExpression(id string) (*entities.Expression, error) {
	var expr entities.Expression
	err := s.db.QueryRow("SELECT id, expression, status, result FROM expressions WHERE id = ?", id).
		Scan(&expr.ID, &expr.Expression, &expr.Status, &expr.Result)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("expression not found")
	}
	return &expr, err
}

func (s *Storage) GetExpressions() ([]entities.Expression, error) {
	rows, err := s.db.Query("SELECT id, expression, status, result FROM expressions ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var expressions []entities.Expression
	for rows.Next() {
		var expr entities.Expression
		err := rows.Scan(&expr.ID, &expr.Expression, &expr.Status, &expr.Result)
		if err != nil {
			return nil, err
		}
		expressions = append(expressions, expr)
	}
	return expressions, nil
}

func (s *Storage) UpdateExpression(id string, status entities.ExpressionStatus, result float64) error {
	_, err := s.db.Exec("UPDATE expressions SET status = ?, result = ? WHERE id = ?",
		status, result, id)
	return err
}
