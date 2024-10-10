package models

type User struct {
	Email        string
	ID           string
	PasswordHash []byte
}

func (user User) GetID() string {
	return user.ID
}

func (user User) GetEmail() string {
	return user.ID
}
