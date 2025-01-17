package domain

import "time"

type TaskStatus string

type Task struct {
	ID          string     `json:"id,omitempty"`
	Title       string     `json:"title" validate:"required"`
	Description string     `json:"description" db:"a_description" validate:"required"`
	Status      TaskStatus `json:"status" db:"a_status" validate:"oneof=pending completed"`
	DueDate     time.Time  `json:"due_date" db:"due_date" validate:"required"`
}

type (
	TaskMutationResponse struct {
		ApiResponse

		Task *Task `json:"task,omitempty"`
	}

	Pagination struct {
		CurrentPage int `json:"current_page"`
		TotalPage   int `json:"total_pages"`
		TotalTask   int `json:"total_task"`
	}

	TaskPaginatedResponse struct {
		Tasks      []Task     `json:"tasks"`
		Pagination Pagination `json:"pagination"`
	}
)
