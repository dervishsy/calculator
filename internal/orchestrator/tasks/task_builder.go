package tasks

import (
	"calculator/internal/orchestrator/parser"
	"calculator/pkg/uuid"
)

func TreeToTasks(root *parser.Node, ExprID string) []Task {
	if root == nil {
		return []Task{}
	}
	tasks := []Task{}

	appendTask(ExprID, root, &tasks)

	// root element to first
	tasks[0], tasks[len(tasks)-1] = tasks[len(tasks)-1], tasks[0]
	return tasks
}

func appendTask(exprID string, root *parser.Node, tasks *[]Task) *Task {
	if root == nil || root.Left == nil || root.Right == nil {
		return nil
	}

	leftArg, rightArg := buildArguments(exprID, root, tasks)

	task := &Task{
		ID:        uuid.New(),
		ExprID:    exprID,
		ArgLeft:   leftArg,
		ArgRight:  rightArg,
		Operation: root.Token.Value,
	}
	*tasks = append(*tasks, *task)
	return task
}

func buildArguments(exprID string, root *parser.Node, tasks *[]Task) (Arg, Arg) {
	var leftArg, rightArg Arg
	if root.Left.Token.Type == parser.Number {
		leftArg = Arg{ArgFloat: root.Left.Value, ArgType: isNumber}
	} else {
		leftArg = Arg{ArgTask: appendTask(exprID, root.Left, tasks), ArgType: isTask}
	}

	if root.Right.Token.Type == parser.Number {
		rightArg = Arg{ArgFloat: root.Right.Value, ArgType: isNumber}
	} else {
		rightArg = Arg{ArgTask: appendTask(exprID, root.Right, tasks), ArgType: isTask}
	}
	return leftArg, rightArg
}
