package postgresql

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/muhrifqii/tuskar/domain"
	"github.com/muhrifqii/tuskar/internal/sqler"
	"go.uber.org/zap"
)

type (
	TaskRepository struct {
		db  *sqler.SqlxWrapper
		log *zap.Logger
	}
)

func NewTaskRepository(db *sqlx.DB, zap *zap.Logger) *TaskRepository {
	return &TaskRepository{
		db:  sqler.NewSqlxWrapper(db, zap),
		log: zap,
	}
}

func (r *TaskRepository) GetAllTasks(c context.Context, params *domain.TaskQueryParams) ([]domain.Task, int, error) {
	tasks := []domain.Task{}

	baseQuery := "SELECT * FROM tasks"
	countQuery := "SELECT COUNT(id) FROM tasks"
	filterQuery := " WHERE 1=1"
	limit := 10
	offset := 0
	argCount := 1
	args := []interface{}{}

	if params.Limit > 0 {
		limit = params.Limit
	}
	if params.Page > 0 {
		offset = limit * (params.Page - 1)
	}
	if params.Status != "" {
		filterQuery += fmt.Sprintf(" AND status = $%d", argCount)
		args = append(args, params.Status)
		argCount++
	}
	if params.Search != "" {
		filterQuery += fmt.Sprintf(" AND (title ILIKE $%d OR a_description ILIKE $%d)", argCount, argCount+1)
		args = append(args, "%"+params.Search+"%", "%"+params.Search+"%")
		argCount += 2
	}

	filterArgs := args[:]
	limitQuery := fmt.Sprintf(" ORDER BY id LIMIT $%d OFFSET $%d", argCount, argCount+1)
	args = append(args, limit, offset)

	count := 0
	err := r.db.QueryRow(countQuery+filterQuery, filterArgs...).Scan(&count)
	r.log.Debug("count", zap.String("query", countQuery+filterQuery), zap.Any("args", filterArgs))
	if err != nil {
		return tasks, count, err
	}
	err = r.db.Select(&tasks, baseQuery+filterQuery+limitQuery, args...)
	r.log.Debug("select", zap.String("query", baseQuery+filterQuery+limitQuery), zap.Any("args", args))
	for i := range tasks {
		tasks[i].DueDate = tasks[i].DueDate[:10]
	}
	return tasks, count, err
}

func (r *TaskRepository) GetByID(c context.Context, id string) (domain.Task, error) {
	var task domain.Task
	err := r.db.Get(&task, "SELECT * FROM tasks WHERE identifier = $1", id)
	if err != nil {
		return task, err
	}
	task.DueDate = task.DueDate[:10]
	return task, err
}

func (r *TaskRepository) CreateTask(c context.Context, task *domain.Task) error {
	_, err := r.db.NamedExec("INSERT INTO tasks (title, a_description, a_status, due_date) VALUES (:title, :a_description, :a_status, :due_date)", task)
	return err
}

func (r *TaskRepository) UpdateTask(c context.Context, task *domain.Task) error {
	result, err := r.db.NamedExec("UPDATE tasks SET title = :title, a_description = :a_description, a_status = :a_status, due_date = :due_date WHERE identifier = :identifier", task)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if affected == 0 {
		return domain.ErrNotFound
	}
	return err
}

func (r *TaskRepository) DeleteTask(c context.Context, id string) error {
	result, err := r.db.Exec("DELETE FROM tasks WHERE identifier = $1", id)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if affected == 0 {
		return domain.ErrNotFound
	}
	return err
}
