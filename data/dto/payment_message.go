package dto

type PaymentMessage struct {
	Email          string `json:"email"`
	Amount         int    `json:"amount"`
	Token          string `json:"token"`
	EmailRecipient string `json:"email_recipient"`
	Type           string `json:"type"`
}
