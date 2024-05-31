package scheduler

import (
	"calculator/internal/orchestrator/parser"
	"calculator/internal/shared/entities"
	"calculator/pkg/uuid"
)

// TreeToTasks converts a binary tree representation of an arithmetic expression
// into a list of tasks.
//
// Parameters:
// - root: The root node of the binary tree representing the arithmetic expression.
// - ExprID: The ID of the expression.
//
// Returns:
// - []Task: The list of tasks representing the arithmetic expression.
func TreeToTasks(root *parser.Node, ExprID string) []entities.Task {
	if root == nil {
		return []entities.Task{}
	}
	tasks := []entities.Task{}

	appendTask(ExprID, root, &tasks)

	// root element to first
	if len(tasks) > 0 {
		tasks[0], tasks[len(tasks)-1] = tasks[len(tasks)-1], tasks[0]
	}
	return tasks
}

func appendTask(exprID string, root *parser.Node, tasks *[]entities.Task) *entities.Task {
	if root == nil || root.Left == nil || root.Right == nil {
		return nil
	}

	leftArg, rightArg := buildArguments(exprID, root, tasks)

	task := &entities.Task{
		ID:        uuid.New(),
		ExprID:    exprID,
		ArgLeft:   leftArg,
		ArgRight:  rightArg,
		Operation: root.Token.Value,
	}
	*tasks = append(*tasks, *task)
	return task
}

func buildArguments(exprID string, root *parser.Node, tasks *[]entities.Task) (entities.Arg, entities.Arg) {
	var leftArg, rightArg entities.Arg
	if root.Left.Token.Type == parser.Number {
		leftArg = entities.Arg{ArgFloat: root.Left.Value, ArgType: entities.IsNumber}
	} else {
		leftArg = entities.Arg{ArgTask: appendTask(exprID, root.Left, tasks), ArgType: entities.IsTask}
	}

	if root.Right.Token.Type == parser.Number {
		rightArg = entities.Arg{ArgFloat: root.Right.Value, ArgType: entities.IsNumber}
	} else {
		rightArg = entities.Arg{ArgTask: appendTask(exprID, root.Right, tasks), ArgType: entities.IsTask}
	}
	return leftArg, rightArg
}
