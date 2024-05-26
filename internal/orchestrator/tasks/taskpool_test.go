package tasks

import "testing"

// TestGetTaskToCompute tests the GetTaskToCompute function of the TaskPool struct.
// It tests the function under different scenarios and verifies that the function
// behaves as expected.
func TestGetTaskToCompute(t *testing.T) {
	// Test case 1: No tasks to compute
	// Create an empty task pool and call the function. Expect an error.
	taskPool := &TaskPool{
		tasks:      map[string]*Task{},
		sentTasks:  map[string]bool{},
		taskOwners: map[string]string{},
	}
	_, err := taskPool.GetTaskToCompute()
	if err == nil {
		t.Errorf("Expected error, got nil")
	}

	// Test case 2: Task with both arguments as numbers and not sent before
	// Create a task pool with one task that has both arguments as numbers and
	// not sent before. Call the function. Expect the task to be returned and
	// marked as sent.
	task := Task{
		ID:        "task1",
		ArgLeft:   Arg{ArgType: isNumber, ArgFloat: 1.0},
		ArgRight:  Arg{ArgType: isNumber, ArgFloat: 2.0},
		Operation: "+",
	}
	taskPool = &TaskPool{
		tasks: map[string]*Task{
			"task1": &task,
		},
		sentTasks:  map[string]bool{},
		taskOwners: map[string]string{},
	}
	resultTask, err := taskPool.GetTaskToCompute()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if resultTask != task {
		t.Errorf("Expected task %v, got %v", task, resultTask)
	}
	if !taskPool.sentTasks["task1"] {
		t.Errorf("Expected task1 to be marked as sent")
	}

	// Test case 3: Task with one argument as a number and the other as a task
	// Create a task pool with one task that has one argument as a number and
	// the other as a task. Also create the task that is referenced by the task.
	// Call the function. Expect an error.
	task = Task{
		ID:        "task2",
		ArgLeft:   Arg{ArgType: isNumber, ArgFloat: 1.0},
		ArgRight:  Arg{ArgType: isTask, ArgTask: &Task{ID: "task3", ArgLeft: Arg{ArgType: isNumber, ArgFloat: 2.0}}},
		Operation: "+",
	}
	taskPool = &TaskPool{
		tasks: map[string]*Task{
			"task2": &task,
			"task3": {ID: "task3", ArgLeft: Arg{ArgType: isNumber, ArgFloat: 2.0}},
		},
		sentTasks:  map[string]bool{},
		taskOwners: map[string]string{"task3": "task2"},
	}
	_, err = taskPool.GetTaskToCompute()
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}
