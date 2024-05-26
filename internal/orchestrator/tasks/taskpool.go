package tasks

import (
	"fmt"
	"sync"
)

type TaskPool struct {
	tasks           map[string]*Task
	taskOwners      map[string]string
	sentTasks       map[string]bool
	expressionsRoot map[string]string
	mu              sync.RWMutex
}

func NewTaskPool() *TaskPool {
	taskPool := &TaskPool{
		tasks:           make(map[string]*Task),
		sentTasks:       make(map[string]bool),
		taskOwners:      make(map[string]string),
		expressionsRoot: make(map[string]string),
		mu:              sync.RWMutex{},
	}

	return taskPool
}

func (tp *TaskPool) AddTasks(tasks []Task) {
	tp.mu.Lock()
	defer tp.mu.Unlock()

	for _, task := range tasks {
		tp.tasks[task.ID] = &task
		if task.ArgLeft.ArgType == isTask {
			tp.taskOwners[task.ArgLeft.ArgTask.ID] = task.ID
		}

		if task.ArgRight.ArgType == isTask {
			tp.taskOwners[task.ArgRight.ArgTask.ID] = task.ID
		}
	}
	tp.expressionsRoot[tasks[0].ID] = tasks[0].ExprID

}

func (tp *TaskPool) GetTaskToCompute() (Task, error) {
	tp.mu.RLock()
	defer tp.mu.RUnlock()
	for _, task := range tp.tasks {
		if task.ArgLeft.ArgType == isNumber &&
			task.ArgRight.ArgType == isNumber &&
			!tp.sentTasks[task.ID] {

			tp.sentTasks[task.ID] = true
			return *task, nil
		}
	}
	return Task{}, fmt.Errorf("no tasks to compute")
}

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
		owner.ArgLeft.ArgType = isNumber
		owner.ArgLeft.ArgFloat = result
	} else if tp.isIdArg(id, owner.ArgRight) {
		owner.ArgRight.ArgType = isNumber
		owner.ArgRight.ArgFloat = result
	} else {
		return fmt.Errorf("task %s owner not found", id)
	}

	return nil
}
func (tp *TaskPool) DeleteTask(id string) {

	tp.mu.Lock()
	defer tp.mu.Unlock()

	delete(tp.sentTasks, id)
	delete(tp.taskOwners, id)
	delete(tp.tasks, id)

}

func (tp *TaskPool) DeleteExpression(id string) {

	tp.mu.Lock()
	defer tp.mu.Unlock()

	delete(tp.expressionsRoot, id)

}

func (tp *TaskPool) IsLastTask(id string) (bool, error) {

	tp.mu.RLock()
	defer tp.mu.RUnlock()

	if _, ok := tp.expressionsRoot[id]; ok {
		return true, nil
	}

	return false, nil
}

func (tp *TaskPool) isIdArg(id string, arg Arg) bool {

	if arg.ArgType == isTask && arg.ArgTask.ID == id {
		return true
	}
	return false
}
