package task

type CreateTaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Project     string `json:"project"`
	Status      string `json:"status"`
	UserID      int    `json:"user_id"`
	Date        string `json:"date"` // Expecting YYYY-MM-DD
	EffortTime  int    `json:"effort_time"`
}

type UpdateTaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Project     string `json:"project"`
	Status      string `json:"status"`
	UserID      int    `json:"user_id"`
	Date        string `json:"date"`
	EffortTime  int    `json:"effort_time"`
}
