package entity

import "time"

// TaskStatus represents the status of a task.
type TaskStatus string

const (
	TaskStatusPending    TaskStatus = "pending"
	TaskStatusInProgress TaskStatus = "in_progress"
	TaskStatusCompleted  TaskStatus = "completed"
)

// Task represents a task item.
type Task struct {
	ID         int        `json:"id" db:"id"`
	UserID     int        `json:"user_id" db:"user_id"`
	Title      string     `json:"title" db:"title"`
	Description string    `json:"description" db:"description"`
	Status     TaskStatus `json:"status" db:"status"`
	Date       time.Time  `json:"date" db:"date"`
	EffortTime int        `json:"effort_time" db:"effort_time"`
	CreatedAt  time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at" db:"updated_at"`
}
