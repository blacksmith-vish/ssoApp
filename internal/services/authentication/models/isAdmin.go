package models

type IsAdminRequest struct {
	UserID string `validate:"required,uuid4"`
}

type IsAdminResponse struct {
	IsAdmin bool
}
