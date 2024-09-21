package models

type RegisterRequest struct {
	Nickname string `json:"nickname" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"pswd" validate:"required"`
}

type RegisterResponse struct {
	UserID string `validate:"required,uuid4"`
}
