package service

import (
	"bytes"
	"email-service/config"
	"email-service/data/dto"
	"fmt"
	"log"
	"net/smtp"
	"text/template"
	"time"
)

type GmailSenderEmail struct {
	cfg *config.Config
}

func NewGmailSenderEmail(cfg *config.Config) SenderEmailService {
	return GmailSenderEmail{cfg: cfg}
}

func (g GmailSenderEmail) SendOrderEmail(message *dto.EmailRequest) error {
	auth := smtp.PlainAuth("", g.cfg.GmailHostConfig.EmailSender,
		g.cfg.GmailHostConfig.Password, g.cfg.GmailHostConfig.StmpHost)

	// Sending email.
	t, err := template.ParseFiles(g.cfg.EmailTemplate.OrderTemplate)
	if err != nil {
		return err
	}

	var body bytes.Buffer

	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body.Write([]byte(fmt.Sprintf("Subject: [Latipe] Xác nhận đơn hàng thành công ! \n%s\n\n", mimeHeaders)))

	t.Execute(&body, struct {
		Name string
		URL  string
		ID   string
	}{
		Name: message.Name,
		URL:  message.Url,
		ID:   message.OrderId,
	})

	err = smtp.SendMail(g.cfg.GmailHostConfig.StmpHost+":"+g.cfg.GmailHostConfig.StmpPort, auth,
		"noreply@latipe.vn", []string{message.EmailRecipient}, body.Bytes())
	if err != nil {
		fmt.Println(err)
		return err
	}
	log.Printf("[Email Order] the email was sent successful at %v", time.Now())
	return nil
}

func (g GmailSenderEmail) SendRegisterEmail(message *dto.EmailRequest) error {
	auth := smtp.PlainAuth("", g.cfg.GmailHostConfig.EmailSender,
		g.cfg.GmailHostConfig.Password, g.cfg.GmailHostConfig.StmpHost)

	// Sending email.
	t, err := template.ParseFiles(g.cfg.EmailTemplate.OrderTemplate)
	if err != nil {
		return err
	}

	var body bytes.Buffer

	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body.Write([]byte(fmt.Sprintf("Subject: [Latipe] Xác nhận đăng ký tài khoản ! \n%s\n\n", mimeHeaders)))

	t.Execute(&body, struct {
		Title   string
		URL     string
		Message string
	}{
		Title:   "Xác nhận email đăng ký tài khoản",
		URL:     message.Url,
		Message: "Hoàn thành xác nhận đăng ký tài khoản để sử dụng các chức năng của hệ thống",
	})

	err = smtp.SendMail(g.cfg.GmailHostConfig.StmpHost+":"+g.cfg.GmailHostConfig.StmpPort, auth,
		"noreply@latipe.vn", []string{message.EmailRecipient}, body.Bytes())
	if err != nil {
		fmt.Println(err)
		return err
	}
	log.Printf("[Email Order] the email was sent successful at %v", time.Now())
	return nil
}

func (g GmailSenderEmail) SendForgotPassword(message *dto.EmailRequest) error {
	auth := smtp.PlainAuth("", g.cfg.GmailHostConfig.EmailSender,
		g.cfg.GmailHostConfig.Password, g.cfg.GmailHostConfig.StmpHost)

	// Sending email.
	t, err := template.ParseFiles(g.cfg.EmailTemplate.OrderTemplate)
	if err != nil {
		return err
	}

	var body bytes.Buffer

	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body.Write([]byte(fmt.Sprintf("Subject: [Latipe] Xác nhận yêu cầu quên mật khẩu ! \n%s\n\n", mimeHeaders)))

	t.Execute(&body, struct {
		Title   string
		URL     string
		Message string
	}{
		Title:   "Đặt lại mật khẩu",
		URL:     message.Url,
		Message: "Click vào link bên dưới để tạo mật khẩu mới cho tài khoản",
	})

	err = smtp.SendMail(g.cfg.GmailHostConfig.StmpHost+":"+g.cfg.GmailHostConfig.StmpPort, auth,
		"noreply@latipe.vn", []string{message.EmailRecipient}, body.Bytes())
	if err != nil {
		fmt.Println(err)
		return err
	}
	log.Printf("[Email Order] the email was sent successful at %v", time.Now())
	return nil
}

func (g GmailSenderEmail) ConfirmLinkEmail(message *dto.EmailRequest) error {
	//TODO implement me
	panic("implement me")
}
