package scheduler

import (
	"calculator/configs"
	"calculator/internal/orchestrator/parser"
	"calculator/internal/orchestrator/storage"
	"calculator/internal/orchestrator/tasks"
	"calculator/internal/shared"
	"calculator/pkg/logger"
	"errors"
	"time"
)

var (
	ErrNoTasksAvailable = errors.New("no tasks available")
	ErrTaskNotFound     = errors.New("task not found")
)

// Scheduler is responsible for managing the execution of arithmetic expressions.
type Scheduler struct {
	cfg      *configs.Config
	storage  *storage.Storage
	taskPoll *tasks.TaskPool
}

// NewScheduler creates a new instance of the Scheduler.
func NewScheduler(storage *storage.Storage, cfg *configs.Config) *Scheduler {
	return &Scheduler{
		cfg:      cfg,
		storage:  storage,
		taskPoll: tasks.NewTaskPool(),
	}
}

// ScheduleExpression schedules an arithmetic expression for execution.
func (s *Scheduler) ScheduleExpression(id, expr string) error {
	if _, err := s.storage.GetExpression(id); err == nil {
		return errors.New("expression already exists")
	}

	rootNode, err := parser.Parse(expr)
	if err != nil {
		return err
	}
	tasksList := tasks.TreeToTasks(rootNode, id)

	s.taskPoll.AddTasks(tasksList)

	return s.storage.CreateExpression(id, expr, tasksList)
}

// GetTask retrieves the next task from the queue.
func (s *Scheduler) GetTask() (*shared.Task, error) {
	task, err := s.taskPoll.GetTaskToCompute()

	if err != nil {
		return nil, ErrNoTasksAvailable
	}
	agentTask := s.taskToAgentTask(task)
	return &agentTask, nil
}

// ProcessResult processes the result of a task computation.
// Deletes the task from the queue after processing.
func (s *Scheduler) ProcessResult(taskID string, result float64) error {

	expr, err := s.storage.GetExpressionByTaskID(taskID)
	if err != nil {
		if err == storage.ErrExpressionNotFound {
			return ErrTaskNotFound
		}
		return err
	}

	err = s.taskPoll.SetTaskResultAfterCompute(taskID, result)
	s.taskPoll.DeleteTask(taskID)
	if err != nil {
		return err
	}

	isLastTask, err := s.taskPoll.IsLastTask(taskID)
	if err != nil {
		return err
	}

	if !isLastTask {
		return nil
	}
	s.taskPoll.DeleteExpression(expr.ID)
	if err = s.storage.UpdateExpression(expr.ID, shared.ExpressionStatusCompleted, result); err != nil {
		return err
	}

	if expr.Status == shared.ExpressionStatusCompleted {
		logger.Infof("Expression %s completed with result %f", expr.ID, result)
	}

	return nil
}

func (s *Scheduler) taskToAgentTask(task tasks.Task) shared.Task {

	return shared.Task{
		ExprID:        task.ExprID,
		ID:            task.ID,
		Arg1:          task.ArgLeft.ArgFloat,
		Arg2:          task.ArgRight.ArgFloat,
		Operation:     task.Operation,
		OperationTime: s.getOperationTime(task.Operation),
	}
}
func (s *Scheduler) getOperationTime(operation string) time.Duration {
	opTime := 0
	switch operation {
	case "+":
		opTime = s.cfg.TimeAdditionMS
	case "-":
		opTime = s.cfg.TimeSubtractionMS
	case "*":
		opTime = s.cfg.TimeMultiplicationMS
	case "/":
		opTime = s.cfg.TimeDivisionMS
	}

	return time.Duration(opTime) * time.Millisecond
}
