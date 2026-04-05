package models

type Todo struct {
	TodoID      int    `json:"todo_id"`
	UserID      int    `json:"user_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
	Priority    string `json:"priority"`
	DueDate     string `json:"due_date"`
	DueTime     string `json:"due_time"`
}
