package models

import "time"

type TaskStatus int

const (
	Pending TaskStatus = iota
	InProgress
	Completed
	Archived
)

type Task struct {
	ID          uint64     `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      TaskStatus `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}
