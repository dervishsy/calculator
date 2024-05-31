package memory_task_storage

import (
	"calculator/internal/shared/entities"
	"fmt"
	"sync"
)

// TaskPool is a struct that represents a task pool in the orchestrator.
type TaskPool struct {
	tasks           map[string]*entities.Task
	taskOwners      map[string]string
	sentTasks       map[string]bool
	expressionsRoot map[string]string
	mu              sync.RWMutex
}

// NewTaskPool creates a new instance of the TaskPool struct.
func NewTaskPool() *TaskPool {
	taskPool := &TaskPool{
		tasks:           make(map[string]*entities.Task),
		sentTasks:       make(map[string]bool),
		taskOwners:      make(map[string]string),
		expressionsRoot: make(map[string]string),
		mu:              sync.RWMutex{},
	}

	return taskPool
}

// AddTasks adds a slice of tasks to the task pool.
func (tp *TaskPool) AddTasks(tasks []entities.Task) {
	tp.mu.Lock()
	defer tp.mu.Unlock()

	for _, task := range tasks {
		tp.tasks[task.ID] = &task
		if task.ArgLeft.ArgType == entities.IsTask {
			tp.taskOwners[task.ArgLeft.ArgTask.ID] = task.ID
		}

		if task.ArgRight.ArgType == entities.IsTask {
			tp.taskOwners[task.ArgRight.ArgTask.ID] = task.ID
		}
	}
	tp.expressionsRoot[tasks[0].ID] = tasks[0].ExprID

}

// GetTaskToCompute returns the next task to compute in the task pool.
func (tp *TaskPool) GetTaskToCompute() (entities.Task, error) {
	tp.mu.RLock()
	defer tp.mu.RUnlock()
	for _, task := range tp.tasks {
		if task.ArgLeft.ArgType == entities.IsNumber &&
			task.ArgRight.ArgType == entities.IsNumber &&
			!tp.sentTasks[task.ID] {

			tp.sentTasks[task.ID] = true
			return *task, nil
		}
	}
	return entities.Task{}, fmt.Errorf("no tasks to compute")
}

// SetTaskResultAfterCompute sets the result of a task after it has been computed.
func (tp *TaskPool) SetTaskResultAfterCompute(id string, result float64) error {

	tp.mu.Lock()
	defer tp.mu.Unlock()

	if _, ok := tp.tasks[id]; !ok {
		return fmt.Errorf("task %s not found", id)
	}

	if _, ok := tp.expressionsRoot[id]; ok {
		return nil
	}

	ownerID, ok := tp.taskOwners[id]

	if !ok {
		return fmt.Errorf("task %s owner not found", id)
	}

	owner := tp.tasks[ownerID]

	if tp.isIdArg(id, owner.ArgLeft) {
		owner.ArgLeft.ArgType = entities.IsNumber
		owner.ArgLeft.ArgFloat = result
	} else if tp.isIdArg(id, owner.ArgRight) {
		owner.ArgRight.ArgType = entities.IsNumber
		owner.ArgRight.ArgFloat = result
	} else {
		return fmt.Errorf("task %s owner not found", id)
	}

	return nil
}

// DeleteTask deletes a task from the task pool.
func (tp *TaskPool) DeleteTask(id string) {

	tp.mu.Lock()
	defer tp.mu.Unlock()

	delete(tp.sentTasks, id)
	delete(tp.taskOwners, id)
	delete(tp.tasks, id)

}

// DeleteExpression deletes an expression from the task pool.
func (tp *TaskPool) DeleteExpression(id string) {

	tp.mu.Lock()
	defer tp.mu.Unlock()

	delete(tp.expressionsRoot, id)

}

// IsLastTask returns true if the task is the last task in the expression.
func (tp *TaskPool) IsLastTask(id string) (bool, error) {

	tp.mu.RLock()
	defer tp.mu.RUnlock()

	if _, ok := tp.expressionsRoot[id]; ok {
		return true, nil
	}

	return false, nil
}

func (tp *TaskPool) GetExpressionIDByTaskID(taskID string) (string, error) {
	tp.mu.RLock()
	defer tp.mu.RUnlock()

	task, ok := tp.tasks[taskID]
	if !ok {
		return "", fmt.Errorf("task %s not found", taskID)
	}

	return task.ExprID, nil
}

func (tp *TaskPool) isIdArg(id string, arg entities.Arg) bool {

	if arg.ArgType == entities.IsTask && arg.ArgTask.ID == id {
		return true
	}
	return false
}
