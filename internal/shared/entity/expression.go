package shared

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
