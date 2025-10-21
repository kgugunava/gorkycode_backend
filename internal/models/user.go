package models

type User struct {
	UserId int
	Name string
	Email string
	PasswordToken string
	Avatar []byte
}