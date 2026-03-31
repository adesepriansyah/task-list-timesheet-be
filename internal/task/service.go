package task

import (
	"context"
	"errors"
	"time"

	"github.com/adesepriansyah/task-list-timesheet-be/internal/entity"
	"github.com/adesepriansyah/task-list-timesheet-be/internal/repository"
)

// Service provides task management logic.
type Service interface {
	CreateTask(ctx context.Context, req CreateTaskRequest) error
	GetTasks(ctx context.Context, filter repository.TaskFilter) ([]entity.Task, error)
	UpdateTask(ctx context.Context, id int, userID int, req UpdateTaskRequest) error
	DeleteTask(ctx context.Context, id int) error
}

type service struct {
	repo repository.TaskRepository
}

// NewService creates a new task service.
func NewService(repo repository.TaskRepository) Service {
	return &service{repo}
}

func (s *service) CreateTask(ctx context.Context, req CreateTaskRequest) error {
	if req.Title == "" {
		return errors.New("title is required")
	}

	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return errors.New("invalid date format, expected YYYY-MM-DD")
	}

	task := &entity.Task{
		UserID:      req.UserID,
		Title:       req.Title,
		Description: req.Description,
		Status:      entity.TaskStatus(req.Status),
		Date:        date,
		EffortTime:  req.EffortTime,
	}

	return s.repo.Create(ctx, task)
}

func (s *service) GetTasks(ctx context.Context, filter repository.TaskFilter) ([]entity.Task, error) {
	if filter.UserID == 0 {
		return nil, errors.New("user_id is required")
	}
	return s.repo.FindAll(ctx, filter)
}

func (s *service) UpdateTask(ctx context.Context, id int, userID int, req UpdateTaskRequest) error {
	existing, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return errors.New("task not found")
	}

	if existing.UserID != userID {
		return errors.New("unauthorized")
	}

	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return errors.New("invalid date format, expected YYYY-MM-DD")
	}

	existing.Title = req.Title
	existing.Description = req.Description
	existing.Status = entity.TaskStatus(req.Status)
	existing.Date = date
	existing.EffortTime = req.EffortTime

	return s.repo.Update(ctx, existing)
}

func (s *service) DeleteTask(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}
