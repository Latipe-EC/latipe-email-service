package dto

type DeliveryAccountMessage struct {
	EmailRecipient string `json:"email_recipient" validator:"email"`
	Email          string `json:"email"`
	Password       string `json:"password"`
}
