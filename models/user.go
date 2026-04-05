package models

type RegisterInput struct {
	Name         string `json:"name"`
	Email        string `json:"email"`
	Password     string `json:"password"`
	ProfileImage string `json:"profile_image"`
}

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type GoogleUser struct {
	ID      string `json:"id"`
	Email   string `json:"email"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}
