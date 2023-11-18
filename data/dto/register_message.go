package dto

type UserRegisterMessage struct {
	EmailRecipient string `json:"email_recipient" validator:"email"`
	Name           string `json:"name"`
	Url            string `json:"url"`
	Code           string `json:"code"`
}
