package dto

type UserRegisterMessage struct {
	Email string `json:"email" validator:"email"`
	Name  string `json:"name"`
	Token string `json:"token"`
}
