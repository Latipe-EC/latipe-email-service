package service

import "email-service/data/dto"

type SenderEmailService interface {
	SendOrderEmail(message *dto.EmailRequest) error
	SendRegisterEmail(message *dto.EmailRequest) error
	SendForgotPassword(message *dto.EmailRequest) error
}
