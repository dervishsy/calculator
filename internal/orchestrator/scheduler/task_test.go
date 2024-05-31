package scheduler

import (
	"calculator/internal/orchestrator/parser"
	"calculator/internal/shared/entities"
	"calculator/pkg/uuid"
	"testing"
)

func TestTreeToTasks(t *testing.T) {
	// Test with nil root
	tasks := TreeToTasks(nil, "TestExprID")
	if len(tasks) != 0 {
		t.Errorf("Expected empty tasks, got %v", tasks)
	}

	// Test with leaf node
	leaf := &parser.Node{
		Token: parser.Token{Type: parser.Number, Value: "2"},
		Value: 2.0,
	}
	tasks = TreeToTasks(leaf, "TestExprID")
	if len(tasks) != 0 {
		t.Errorf("Expected empty tasks, got %v", tasks)
	}

	// Test with single operation node
	singleOp := &parser.Node{
		Token: parser.Token{Type: parser.Plus, Value: "+"},
		Left:  &parser.Node{Token: parser.Token{Type: parser.Number, Value: "2"}, Value: 2.0},
		Right: &parser.Node{Token: parser.Token{Type: parser.Number, Value: "3"}, Value: 3.0},
	}
	tasks = TreeToTasks(singleOp, "TestExprID")
	if len(tasks) != 1 {
		t.Errorf("Expected 1 task, got %v", tasks)
	}
	expectedTask := entities.Task{
		ID:        uuid.New(),
		ExprID:    "TestExprID",
		ArgLeft:   entities.Arg{ArgFloat: 2.0, ArgType: entities.IsNumber},
		ArgRight:  entities.Arg{ArgFloat: 3.0, ArgType: entities.IsNumber},
		Operation: "+",
	}
	if !Compare(tasks[0], expectedTask) {
		t.Errorf("Expected task %v, got %v", expectedTask, tasks[0])
	}

	// Test with nested operation node
	nestedOp := &parser.Node{
		Token: parser.Token{Type: parser.Plus, Value: "+"},
		Left: &parser.Node{
			Token: parser.Token{Type: parser.Multiply, Value: "*"},
			Left:  &parser.Node{Token: parser.Token{Type: parser.Number, Value: "2"}, Value: 2.0},
			Right: &parser.Node{Token: parser.Token{Type: parser.Number, Value: "3"}, Value: 3.0},
		},
		Right: &parser.Node{Token: parser.Token{Type: parser.Number, Value: "4"}, Value: 4.0},
	}
	tasks = TreeToTasks(nestedOp, "TestExprID")
	if len(tasks) != 2 {
		t.Errorf("Expected 2 tasks, got %v", tasks)
	}
	expectedTask = entities.Task{
		ID:        uuid.New(),
		ExprID:    "TestExprID",
		ArgLeft:   entities.Arg{ArgFloat: 2.0, ArgType: entities.IsNumber},
		ArgRight:  entities.Arg{ArgFloat: 3.0, ArgType: entities.IsNumber},
		Operation: "*",
	}
	if !Compare(tasks[0], expectedTask) {
		t.Errorf("Expected task %v, got %v", expectedTask, tasks[0])
	}
	expectedTask = entities.Task{
		ID:        uuid.New(),
		ExprID:    "TestExprID",
		ArgLeft:   entities.Arg{ArgTask: &tasks[0], ArgType: entities.IsTask},
		ArgRight:  entities.Arg{ArgFloat: 4.0, ArgType: entities.IsNumber},
		Operation: "+",
	}
	if !Compare(tasks[1], expectedTask) {
		t.Errorf("Expected task %v, got %v", expectedTask, tasks[1])
	}
}

// Compare compares two tasks and returns true if they are the same.
func Compare(t, t2 entities.Task) bool {
	return t.ExprID == t2.ExprID &&
		t.Operation == t2.Operation &&
		t.ArgLeft.ArgFloat == t2.ArgLeft.ArgFloat &&
		t.ArgRight.ArgFloat == t2.ArgRight.ArgFloat
}
