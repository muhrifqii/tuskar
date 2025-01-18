package domain

type TaskStatus string

type Task struct {
	DbID        int64      `json:"-" db:"id"`
	ID          string     `json:"id,omitempty" db:"identifier"`
	Title       string     `json:"title" validate:"required"`
	Description string     `json:"description" db:"a_description" validate:"required"`
	Status      TaskStatus `json:"status" db:"a_status" validate:"required,oneof=pending completed"`
	DueDate     string     `json:"due_date" db:"due_date" validate:"required,datetime=2006-01-02"`
}

type (
	TaskMutationResponse struct {
		Message string `json:"message"`
		Task    *Task  `json:"task,omitempty"`
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

	TaskQueryParams struct {
		Status TaskStatus `query:"status"`
		Page   int        `query:"page"`
		Limit  int        `query:"limit"`
		Search string     `query:"search"`
	}
)
