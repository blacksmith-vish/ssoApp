package models

type RegisterRequest struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required"`
}

type RegisterResponse struct {
	UserID string `validate:"required,uuid4"`
}
