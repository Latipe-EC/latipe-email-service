package dto

const (
	ORDER    = 1
	REGISTER = 2
	FORGOT   = 3
)

type EmailRequest struct {
	EmailRecipient string `json:"email_recipient" validator:"email"`
	Name           string `json:"name"`
	OrderId        string `json:"order_id"`
	Url            string `json:"url"`
	Code           string `json:"code"`
	Type           int    `json:"type"`
}
