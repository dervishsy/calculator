package scheduler

import (
	"calculator/configs"
	"calculator/internal/orchestrator/parser"
	"calculator/internal/orchestrator/tasks"
	"calculator/internal/shared/entities"
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
	storage  ExpressionService
	taskPoll TaskService
}

// NewScheduler creates a new instance of the Scheduler.
func NewScheduler(storage ExpressionService, cfg *configs.Config) *Scheduler {
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

	return s.storage.CreateExpression(id, expr)
}

// GetTask retrieves the next task from the queue.
func (s *Scheduler) GetTask() (*entities.AgentTask, error) {
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

	exprID, err := s.taskPoll.GetExpressionIDByTaskID(taskID)

	if err != nil {
		return ErrTaskNotFound
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

	if isLastTask {
		logger.Infof("Expression %s completed with result %f", exprID, result)
	}

	s.taskPoll.DeleteExpression(exprID)
	if err = s.storage.UpdateExpression(exprID, entities.ExpressionStatusCompleted, result); err != nil {
		return err
	}

	return nil
}

func (s *Scheduler) taskToAgentTask(task entities.Task) entities.AgentTask {

	return entities.AgentTask{
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
