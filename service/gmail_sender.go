package service

import (
	"bytes"
	"email-service/config"
	"email-service/data/dto"
	"fmt"
	"log"
	"net/smtp"
	"strings"
	"text/template"
	"time"
)

type GmailSenderEmail struct {
	cfg *config.Config
}

func NewGmailSenderEmail(cfg *config.Config) SenderEmailService {
	return GmailSenderEmail{cfg: cfg}
}

func (g GmailSenderEmail) SendOrderEmail(message *dto.OrderMessage) error {
	auth := smtp.PlainAuth("", g.cfg.GmailHostConfig.EmailSender,
		g.cfg.GmailHostConfig.Password, g.cfg.GmailHostConfig.StmpHost)

	confirmUrl := fmt.Sprintf("%s/order/%s", g.cfg.HostURL, message.OrderId)
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
		URL:  confirmUrl,
		ID:   strings.ToUpper(message.OrderId),
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

func (g GmailSenderEmail) SendRegisterEmail(message *dto.UserRegisterMessage) error {
	auth := smtp.PlainAuth("", g.cfg.GmailHostConfig.EmailSender,
		g.cfg.GmailHostConfig.Password, g.cfg.GmailHostConfig.StmpHost)

	// Sending email.
	t, err := template.ParseFiles(g.cfg.EmailTemplate.ConfirmLinkTemplate)
	if err != nil {
		return err
	}
	confirmUrl := fmt.Sprintf("%s/auth/verify-account/%s", g.cfg.HostURL, message.Token)

	var body bytes.Buffer

	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body.Write([]byte(fmt.Sprintf("Subject: [Latipe] Xác nhận đăng ký tài khoản ! \n%s\n\n", mimeHeaders)))

	t.Execute(&body, struct {
		Title   string
		Url     string
		Message string
	}{
		Title:   "Xác nhận email đăng ký tài khoản",
		Url:     confirmUrl,
		Message: "Hoàn thành xác nhận đăng ký tài khoản để sử dụng các chức năng của hệ thống",
	})

	err = smtp.SendMail(g.cfg.GmailHostConfig.StmpHost+":"+g.cfg.GmailHostConfig.StmpPort, auth,
		"noreply@latipe.vn", []string{message.Email}, body.Bytes())
	if err != nil {
		fmt.Println(err)
		return err
	}
	log.Printf("[Email Order] the email was sent successful at %v", time.Now())
	return nil
}

func (g GmailSenderEmail) SendForgotPassword(message *dto.ForgotPasswordMessage) error {
	auth := smtp.PlainAuth("", g.cfg.GmailHostConfig.EmailSender,
		g.cfg.GmailHostConfig.Password, g.cfg.GmailHostConfig.StmpHost)

	// Sending email.
	t, err := template.ParseFiles(g.cfg.EmailTemplate.ConfirmLinkTemplate)
	if err != nil {
		return err
	}

	var body bytes.Buffer

	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body.Write([]byte(fmt.Sprintf("Subject: [Latipe] Xác nhận yêu cầu quên mật khẩu ! \n%s\n\n", mimeHeaders)))

	confirmUrl := fmt.Sprintf("%s/auth/reset-password/%s", g.cfg.HostURL, message.Token)

	t.Execute(&body, struct {
		Title   string
		Url     string
		Message string
	}{
		Title:   "Đặt lại mật khẩu",
		Url:     confirmUrl,
		Message: "Click vào link bên dưới để tạo mật khẩu mới cho tài khoản",
	})

	err = smtp.SendMail(g.cfg.GmailHostConfig.StmpHost+":"+g.cfg.GmailHostConfig.StmpPort, auth,
		"noreply@latipe.vn", []string{message.Email}, body.Bytes())
	if err != nil {
		fmt.Println(err)
		return err
	}
	log.Printf("[Email Order] the email was sent successful at %v", time.Now())
	return nil
}

func (g GmailSenderEmail) SendDeliveryAccount(message *dto.DeliveryAccountMessage) error {
	auth := smtp.PlainAuth("", g.cfg.GmailHostConfig.EmailSender,
		g.cfg.GmailHostConfig.Password, g.cfg.GmailHostConfig.StmpHost)

	// Sending email.
	t, err := template.ParseFiles(g.cfg.EmailTemplate.DeliveryAccountTemplate)
	if err != nil {
		return err
	}

	url := g.cfg.HostURL + "/auth/active?token=" + message.Token
	var body bytes.Buffer

	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body.Write([]byte(fmt.Sprintf("Subject: [Latipe] Thông báo tạo tài khoản đối tác ! \n%s\n\n", mimeHeaders)))

	t.Execute(&body, struct {
		Email    string
		Password string
		Url      string
	}{
		Email:    message.Email,
		Password: message.Password,
		Url:      url,
	})

	err = smtp.SendMail(g.cfg.GmailHostConfig.StmpHost+":"+g.cfg.GmailHostConfig.StmpPort, auth,
		"noreply@latipe.vn", []string{message.Email}, body.Bytes())
	if err != nil {
		fmt.Println(err)
		return err
	}
	log.Printf("[Delivery email] the email was sent successful at %v", time.Now())
	return nil
}

func (g GmailSenderEmail) SendTakeoutPayment(message *dto.PaymentMessage) error {
	auth := smtp.PlainAuth("", g.cfg.GmailHostConfig.EmailSender,
		g.cfg.GmailHostConfig.Password, g.cfg.GmailHostConfig.StmpHost)

	// Sending email.
	t, err := template.ParseFiles(g.cfg.EmailTemplate.ConfirmTakeoutTemplate)
	if err != nil {
		return err
	}

	url := g.cfg.HostURL + "/withdraw-success/" + message.Token
	var body bytes.Buffer

	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body.Write([]byte(fmt.Sprintf("Subject: [Latipe] Thông báo tạo tài khoản đối tác ! \n%s\n\n", mimeHeaders)))

	t.Execute(&body, struct {
		Email  string
		Type   string
		Amount int
		Url    string
	}{
		Email:  message.Email,
		Type:   message.Type,
		Url:    url,
		Amount: message.Amount,
	})

	err = smtp.SendMail(g.cfg.GmailHostConfig.StmpHost+":"+g.cfg.GmailHostConfig.StmpPort, auth,
		"noreply@latipe.vn", []string{message.Email}, body.Bytes())
	if err != nil {
		fmt.Println(err)
		return err
	}
	log.Printf("[Delivery email] the email was sent successful at %v", time.Now())
	return nil
}
