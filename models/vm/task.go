package vm

import "bitbucket.org/task_service/models"

type Task struct {
	ID          uint64            `json:"id" uri:"id" bindging:"required"`
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Status      models.TaskStatus `json:"status"`
	CreatedAt   int64             `json:"created_at"`
	UpdatedAt   int64             `json:"updated_at"`
}

type CreateTaskRequest struct {
	Title       string            `json:"title" bindging:"required"`
	Description string            `json:"description"`
	Status      models.TaskStatus `json:"status" bindging:"required"`
}

type TaskRequest struct {
	Task
}

type TaskResponse struct {
	Task
}

type ListAllTasksRequest struct {
	PaginationRequest
}

type ListAllTaskResponse struct {
	Tasks []Task `json:"tasks"`
	MetaResponse
}
