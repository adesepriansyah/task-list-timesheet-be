package task

import (
	"context"
	"testing"

	"github.com/adesepriansyah/task-list-timesheet-be/internal/entity"
	"github.com/adesepriansyah/task-list-timesheet-be/internal/repository"
)

type mockTaskRepo struct {
	tasks map[int]*entity.Task
}

func (m *mockTaskRepo) Create(ctx context.Context, task *entity.Task) error {
	task.ID = len(m.tasks) + 1
	m.tasks[task.ID] = task
	return nil
}

func (m *mockTaskRepo) FindByID(ctx context.Context, id int) (*entity.Task, error) {
	return m.tasks[id], nil
}

func (m *mockTaskRepo) FindAll(ctx context.Context, filter repository.TaskFilter) ([]entity.Task, error) {
	var result []entity.Task
	for _, t := range m.tasks {
		if t.UserID == filter.UserID {
			result = append(result, *t)
		}
	}
	return result, nil
}

func (m *mockTaskRepo) Update(ctx context.Context, task *entity.Task) error {
	m.tasks[task.ID] = task
	return nil
}

func (m *mockTaskRepo) Delete(ctx context.Context, id int) error {
	delete(m.tasks, id)
	return nil
}

func TestCreateTask(t *testing.T) {
	repo := &mockTaskRepo{tasks: make(map[int]*entity.Task)}
	svc := NewService(repo)

	t.Run("Success", func(t *testing.T) {
		req := CreateTaskRequest{
			Title:      "Test Task",
			UserID:     1,
			Date:       "2026-04-01",
			Status:     "pending",
			EffortTime: 30,
		}
		err := svc.CreateTask(context.Background(), req)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	t.Run("InvalidDate", func(t *testing.T) {
		req := CreateTaskRequest{
			Title:  "Test Task",
			UserID: 1,
			Date:   "invalid-date",
		}
		err := svc.CreateTask(context.Background(), req)
		if err == nil {
			t.Error("expected error for invalid date format")
		}
	})
}

func TestGetTasks(t *testing.T) {
	repo := &mockTaskRepo{
		tasks: map[int]*entity.Task{
			1: {ID: 1, UserID: 1, Title: "Task 1"},
			2: {ID: 2, UserID: 2, Title: "Task 2"},
		},
	}
	svc := NewService(repo)

	t.Run("Success", func(t *testing.T) {
		tasks, err := svc.GetTasks(context.Background(), repository.TaskFilter{UserID: 1})
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if len(tasks) != 1 {
			t.Errorf("expected 1 task, got %d", len(tasks))
		}
	})

	t.Run("MissingUserID", func(t *testing.T) {
		_, err := svc.GetTasks(context.Background(), repository.TaskFilter{})
		if err == nil {
			t.Error("expected error for missing user_id")
		}
	})
}
