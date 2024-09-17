package models

type LoginRequest struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required"`
	AppID    string `validate:"required,uuid4"`
}

type LoginResponse struct {
	Token string `validate:"required,jwt"`
}
