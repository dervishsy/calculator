package sqlite_task_storage

import (
	"calculator/internal/orchestrator/impl/sqlite"
	"calculator/internal/shared/entities"
	"database/sql"
	"encoding/json"
	"fmt"
)

type TaskPool struct {
	db *sqlite.SQLiteDB
}

func NewTaskPool(db *sqlite.SQLiteDB) *TaskPool {
	return &TaskPool{db: db}
}

func (tp *TaskPool) AddTasks(tasks []entities.Task) error {
	tx, err := tp.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, task := range tasks {
		argLeft, _ := json.Marshal(task.ArgLeft)
		argRight, _ := json.Marshal(task.ArgRight)

		_, err = tx.Exec("INSERT INTO tasks (id, expr_id, arg_left, arg_right, operation) VALUES (?, ?, ?, ?, ?)",
			task.ID, task.ExprID, argLeft, argRight, task.Operation)
		if err != nil {
			return err
		}

		if task.ArgLeft.ArgType == entities.IsTask {
			_, err = tx.Exec("INSERT INTO task_owners (child_id, parent_id) VALUES (?, ?)",
				task.ArgLeft.ArgTask.ID, task.ID)
			if err != nil {
				return err
			}
		}

		if task.ArgRight.ArgType == entities.IsTask {
			_, err = tx.Exec("INSERT INTO task_owners (child_id, parent_id) VALUES (?, ?)",
				task.ArgRight.ArgTask.ID, task.ID)
			if err != nil {
				return err
			}
		}
	}

	_, err = tx.Exec("INSERT INTO expressions_root (task_id, expr_id) VALUES (?, ?)",
		tasks[0].ID, tasks[0].ExprID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (tp *TaskPool) GetTaskToCompute() (entities.Task, error) {
	var task entities.Task
	var argLeftBytes, argRightBytes []byte

	err := tp.db.QueryRow(`
        SELECT id, expr_id, arg_left, arg_right, operation
        FROM tasks
        WHERE id NOT IN (SELECT task_id FROM sent_tasks)
        AND json_extract(arg_left, '$.ArgType') = ?
        AND json_extract(arg_right, '$.ArgType') = ?
        LIMIT 1
    `, entities.IsNumber, entities.IsNumber).Scan(
		&task.ID, &task.ExprID, &argLeftBytes, &argRightBytes, &task.Operation)

	if err != nil {
		return entities.Task{}, fmt.Errorf("no tasks to compute")
	}

	json.Unmarshal(argLeftBytes, &task.ArgLeft)
	json.Unmarshal(argRightBytes, &task.ArgRight)

	_, err = tp.db.Exec("INSERT INTO sent_tasks (task_id) VALUES (?)", task.ID)
	if err != nil {
		return entities.Task{}, err
	}

	return task, nil
}

func (tp *TaskPool) SetTaskResultAfterCompute(id string, result float64) error {
	tx, err := tp.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var ownerID string
	err = tx.QueryRow("SELECT parent_id FROM task_owners WHERE child_id = ?", id).Scan(&ownerID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil // This is the root task
		}
		return err
	}

	var ownerArgLeft, ownerArgRight []byte
	err = tx.QueryRow("SELECT arg_left, arg_right FROM tasks WHERE id = ?", ownerID).Scan(&ownerArgLeft, &ownerArgRight)
	if err != nil {
		return err
	}

	var argLeft, argRight entities.Arg
	json.Unmarshal(ownerArgLeft, &argLeft)
	json.Unmarshal(ownerArgRight, &argRight)

	if argLeft.ArgType == entities.IsTask && argLeft.ArgTask.ID == id {
		argLeft.ArgType = entities.IsNumber
		argLeft.ArgFloat = result
		updatedArgLeft, _ := json.Marshal(argLeft)
		_, err = tx.Exec("UPDATE tasks SET arg_left = ? WHERE id = ?", updatedArgLeft, ownerID)
	} else if argRight.ArgType == entities.IsTask && argRight.ArgTask.ID == id {
		argRight.ArgType = entities.IsNumber
		argRight.ArgFloat = result
		updatedArgRight, _ := json.Marshal(argRight)
		_, err = tx.Exec("UPDATE tasks SET arg_right = ? WHERE id = ?", updatedArgRight, ownerID)
	}

	if err != nil {
		return err
	}

	return tx.Commit()
}

func (tp *TaskPool) DeleteTask(id string) error {
	_, err := tp.db.Exec("DELETE FROM tasks WHERE id = ?", id)
	if err != nil {
		return err
	}
	_, err = tp.db.Exec("DELETE FROM sent_tasks WHERE task_id = ?", id)
	if err != nil {
		return err
	}
	_, err = tp.db.Exec("DELETE FROM task_owners WHERE child_id = ? OR parent_id = ?", id, id)
	return err
}

func (tp *TaskPool) DeleteExpression(id string) error {
	_, err := tp.db.Exec("DELETE FROM expressions_root WHERE expr_id = ?", id)
	return err
}

func (tp *TaskPool) IsLastTask(id string) (bool, error) {
	var count int
	err := tp.db.QueryRow("SELECT COUNT(*) FROM expressions_root WHERE task_id = ?", id).Scan(&count)
	return count > 0, err
}

func (tp *TaskPool) GetExpressionIDByTaskID(taskID string) (string, error) {
	var exprID string
	err := tp.db.QueryRow("SELECT expr_id FROM tasks WHERE id = ?", taskID).Scan(&exprID)
	if err != nil {
		return "", fmt.Errorf("task %s not found", taskID)
	}
	return exprID, nil
}
