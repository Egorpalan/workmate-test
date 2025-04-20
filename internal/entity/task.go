package entity

import (
	"encoding/json"
	"time"
)

type TaskStatus string

const (
	TaskStatusPending    TaskStatus = "pending"
	TaskStatusProcessing TaskStatus = "processing"
	TaskStatusCompleted  TaskStatus = "completed"
	TaskStatusFailed     TaskStatus = "failed"
)

type Task struct {
	ID        string          `json:"id" db:"id"`
	Status    TaskStatus      `json:"status" db:"status"`
	Result    json.RawMessage `json:"result,omitempty" db:"result"`
	Error     string          `json:"error,omitempty" db:"error"`
	CreatedAt time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt time.Time       `json:"updated_at" db:"updated_at"`
}
