package memory_task_storage

import (
	"calculator/internal/shared/entities"
	"testing"
)

// TestGetTaskToCompute tests the GetTaskToCompute function of the TaskPool struct.
// It tests the function under different scenarios and verifies that the function
// behaves as expected.
func TestGetTaskToCompute(t *testing.T) {
	// Test case 1: No tasks to compute
	// Create an empty task pool and call the function. Expect an error.
	taskPool := &TaskPool{
		tasks:      map[string]*entities.Task{},
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
	task := entities.Task{
		ID:        "task1",
		ArgLeft:   entities.Arg{ArgType: entities.IsNumber, ArgFloat: 1.0},
		ArgRight:  entities.Arg{ArgType: entities.IsNumber, ArgFloat: 2.0},
		Operation: "+",
	}
	taskPool = &TaskPool{
		tasks: map[string]*entities.Task{
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
	task = entities.Task{
		ID:        "task2",
		ArgLeft:   entities.Arg{ArgType: entities.IsNumber, ArgFloat: 1.0},
		ArgRight:  entities.Arg{ArgType: entities.IsTask, ArgTask: &entities.Task{ID: "task3", ArgLeft: entities.Arg{ArgType: entities.IsNumber, ArgFloat: 2.0}}},
		Operation: "+",
	}
	taskPool = &TaskPool{
		tasks: map[string]*entities.Task{
			"task2": &task,
			"task3": {ID: "task3", ArgLeft: entities.Arg{ArgType: entities.IsNumber, ArgFloat: 2.0}},
		},
		sentTasks:  map[string]bool{},
		taskOwners: map[string]string{"task3": "task2"},
	}
	_, err = taskPool.GetTaskToCompute()
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}
