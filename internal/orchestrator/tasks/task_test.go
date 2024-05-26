package tasks

import (
	"calculator/internal/orchestrator/parser"
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
	expectedTask := Task{
		ID:        uuid.New(),
		ExprID:    "TestExprID",
		ArgLeft:   Arg{ArgFloat: 2.0, ArgType: isNumber},
		ArgRight:  Arg{ArgFloat: 3.0, ArgType: isNumber},
		Operation: "+",
	}
	if !tasks[0].Compare(expectedTask) {
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
	expectedTask = Task{
		ID:        uuid.New(),
		ExprID:    "TestExprID",
		ArgLeft:   Arg{ArgFloat: 2.0, ArgType: isNumber},
		ArgRight:  Arg{ArgFloat: 3.0, ArgType: isNumber},
		Operation: "*",
	}
	if !tasks[0].Compare(expectedTask) {
		t.Errorf("Expected task %v, got %v", expectedTask, tasks[0])
	}
	expectedTask = Task{
		ID:        uuid.New(),
		ExprID:    "TestExprID",
		ArgLeft:   Arg{ArgTask: &tasks[0], ArgType: isTask},
		ArgRight:  Arg{ArgFloat: 4.0, ArgType: isNumber},
		Operation: "+",
	}
	if !tasks[1].Compare(expectedTask) {
		t.Errorf("Expected task %v, got %v", expectedTask, tasks[1])
	}
}
