package entities

// TaskResult represents the result of a task computation.
type TaskResult struct {
	ID     string  `json:"id"`
	Result float64 `json:"result"`
}
