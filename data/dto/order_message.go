package dto

type OrderMessage struct {
	EmailRecipient string `json:"email_recipient" validator:"email"`
	Name           string `json:"name"`
	OrderId        string `json:"order_id"`
	Code           string `json:"code"`
}
