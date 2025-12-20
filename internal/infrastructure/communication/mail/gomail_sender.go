package mail

import (
	"bytes"
	"context"
	"fmt"
	"hona/backend/bootstrap"
	"html/template"
	"strconv"
	"time"

	"github.com/wneessen/go-mail"
)

type EmailService struct {
	client *mail.Client
}

func NewEmailService() *EmailService {
	config := bootstrap.Run().Env.EmailConfig
	port, _ := strconv.Atoi(config.Port)
	client, err := mail.NewClient(
		config.Host,
		mail.WithPort(port),
		mail.WithSMTPAuth(mail.SMTPAuthPlain),
		mail.WithUsername(config.Username),
		mail.WithPassword(config.Password),
		mail.WithTLSPortPolicy(mail.TLSOpportunistic),
	)

	if err != nil {
		panic(err)
	}

	return &EmailService{
		client: client,
	}
}

func (s *EmailService) SendEmail(to string, subject string, templateFileName string, data interface{}) error {
	m := mail.NewMsg()

	if err := m.From(bootstrap.Run().Env.EmailConfig.From); err != nil {
		return err
	}

	if err := m.To(to); err != nil {
		return err
	}

	m.Subject(subject)

	tmpl, err := template.ParseFiles(bootstrap.Run().Constants.TemplatesPath.Path + templateFileName)
	if err != nil {
		return err
	}
	var body bytes.Buffer
	err = tmpl.Execute(&body, data)
	if err != nil {
		return err
	}
	m.SetBodyString(mail.TypeTextHTML, body.String())

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.client.DialAndSendWithContext(ctx, m); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
