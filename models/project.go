package models

type Project struct {
	ProjectID   int    `json:"project_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CreatedBy   int    `json:"created_by"`
	CreatedAt   string `json:"created_at"`
}

type CreateProjectInput struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}
