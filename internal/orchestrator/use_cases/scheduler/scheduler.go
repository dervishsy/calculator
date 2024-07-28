package scheduler

import (
	use_cases_errors "calculator/internal/orchestrator/use_cases/errors"
	"calculator/internal/orchestrator/use_cases/parser"
	"calculator/internal/shared/configs"
	"calculator/internal/shared/entities"
	"calculator/pkg/logger"
	"time"
)

// Scheduler is responsible for managing the execution of arithmetic expressions.
type Scheduler struct {
	cfg      *configs.Config
	storage  ExpressionService
	taskPoll TaskService
}

// NewScheduler creates a new instance of the Scheduler.
func NewScheduler(storage ExpressionService, task_poll TaskService, cfg *configs.Config) *Scheduler {
	return &Scheduler{
		cfg:      cfg,
		storage:  storage,
		taskPoll: task_poll,
	}
}

// ScheduleExpression schedules an arithmetic expression for execution.
func (s *Scheduler) ScheduleExpression(id, expr string) error {
	if _, err := s.storage.GetExpression(id); err == nil {
		logger.Errorf("Expression with ID %s already exists", id)
		return use_cases_errors.ErrExpressionExists
	}

	rootNode, err := parser.Parse(expr)
	if err != nil {
		return err
	}
	tasksList := TreeToTasks(rootNode, id)

	err = s.taskPoll.AddTasks(tasksList)

	if err != nil {
		logger.Error(err)
		return err
	}

	return s.storage.CreateExpression(id, expr)
}

// GetTask retrieves the next task from the queue.
func (s *Scheduler) GetTask() (*entities.AgentTask, error) {
	task, err := s.taskPoll.GetTaskToCompute()

	if err != nil {
		logger.Error(err)
		return nil, use_cases_errors.ErrNoTasksAvailable
	}
	agentTask := s.taskToAgentTask(task)
	return &agentTask, nil
}

// ProcessResult processes the result of a task computation.
// Deletes the task from the queue after processing.
func (s *Scheduler) ProcessResult(taskID string, result float64) error {

	exprID, err := s.taskPoll.GetExpressionIDByTaskID(taskID)

	if err != nil {
		logger.Error(err)
		return use_cases_errors.ErrNoTasksAvailable
	}

	err = s.taskPoll.SetTaskResultAfterCompute(taskID, result)
	if err != nil {
		logger.Error(err)
		return err
	}
	err = s.taskPoll.DeleteTask(taskID)
	if err != nil {
		logger.Error(err)
		return err
	}

	isLastTask, err := s.taskPoll.IsLastTask(taskID)
	if err != nil {
		logger.Error(err)
		return err
	}

	if !isLastTask {
		return nil
	}

	if isLastTask {
		logger.Infof("Expression %s completed with result %f", exprID, result)
	}

	err = s.taskPoll.DeleteExpression(exprID)
	if err != nil {
		logger.Error(err)
		return err
	}

	if err = s.storage.UpdateExpression(exprID, entities.ExpressionStatusCompleted, result); err != nil {
		logger.Error(err)
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

// GetExpression retrieves an arithmetic expression by its ID.
func (s *Scheduler) GetExpression(id string) (*entities.Expression, error) {
	return s.storage.GetExpression(id)
}

// GetExpressions retrieves all arithmetic expressions.
func (s *Scheduler) GetExpressions() ([]entities.Expression, error) {
	return s.storage.GetExpressions()
}
