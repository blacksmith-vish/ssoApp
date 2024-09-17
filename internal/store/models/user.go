package models

type User struct {
	Email        string
	ID           string
	PasswordHash []byte
}
