package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/adesepriansyah/task-list-timesheet-be/internal/entity"
)

// TaskFilter represents the filter for fetching tasks.
type TaskFilter struct {
	UserID   int
	Search   string
	DateFrom string
	DateTo   string
}

// TaskRepository interface for task operations.
type TaskRepository interface {
	Create(ctx context.Context, task *entity.Task) error
	FindByID(ctx context.Context, id int) (*entity.Task, error)
	FindAll(ctx context.Context, filter TaskFilter) ([]entity.Task, error)
	Update(ctx context.Context, task *entity.Task) error
	Delete(ctx context.Context, id int) error
}

type taskRepository struct {
	db *sql.DB
}

// NewTaskRepository creates a new task repository.
func NewTaskRepository(db *sql.DB) TaskRepository {
	return &taskRepository{db}
}

func (r *taskRepository) Create(ctx context.Context, task *entity.Task) error {
	query := `INSERT INTO tasks (user_id, title, description, project, status, date, effort_time)
              VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id, created_at, updated_at`
	return r.db.QueryRowContext(ctx, query,
		task.UserID, task.Title, task.Description, task.Project, task.Status, task.Date, task.EffortTime,
	).Scan(&task.ID, &task.CreatedAt, &task.UpdatedAt)
}

func (r *taskRepository) FindByID(ctx context.Context, id int) (*entity.Task, error) {
	var task entity.Task
	query := `SELECT id, user_id, title, description, project, status, date, effort_time, created_at, updated_at
              FROM tasks WHERE id = $1`
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&task.ID, &task.UserID, &task.Title, &task.Description, &task.Project, &task.Status, &task.Date, &task.EffortTime, &task.CreatedAt, &task.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *taskRepository) FindAll(ctx context.Context, filter TaskFilter) ([]entity.Task, error) {
	var tasks []entity.Task
	query := `SELECT id, user_id, title, description, project, status, date, effort_time, created_at, updated_at
              FROM tasks WHERE user_id = $1`
	args := []interface{}{filter.UserID}
	argID := 2

	if filter.Search != "" {
		query += fmt.Sprintf(" AND (title ILIKE $%d OR description ILIKE $%d)", argID, argID)
		args = append(args, "%"+filter.Search+"%")
		argID++
	}

	if filter.DateFrom != "" {
		query += fmt.Sprintf(" AND date >= $%d", argID)
		args = append(args, filter.DateFrom)
		argID++
	}

	if filter.DateTo != "" {
		query += fmt.Sprintf(" AND date <= $%d", argID)
		args = append(args, filter.DateTo)
		argID++
	}

	query += " ORDER BY date DESC, created_at DESC"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var task entity.Task
		if err := rows.Scan(
			&task.ID, &task.UserID, &task.Title, &task.Description, &task.Project, &task.Status, &task.Date, &task.EffortTime, &task.CreatedAt, &task.UpdatedAt,
		); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (r *taskRepository) Update(ctx context.Context, task *entity.Task) error {
	query := `UPDATE tasks SET title=$1, description=$2, project=$3, status=$4, date=$5, effort_time=$6, updated_at=NOW()
              WHERE id=$7 AND user_id=$8`
	_, err := r.db.ExecContext(ctx, query,
		task.Title, task.Description, task.Project, task.Status, task.Date, task.EffortTime, task.ID, task.UserID,
	)
	return err
}

func (r *taskRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM tasks WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
