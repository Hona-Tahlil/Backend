package consumers

import (
	"encoding/json"
	"hona/backend/bootstrap"
	domainmail "hona/backend/internal/domain/mail"
	"hona/backend/internal/infrastructure/rabbitmq"
	"log/slog"
)

var constants = bootstrap.Run().Constants.RabbitMQConstants

type EmailConsumer struct {
	rabbitMQ     *rabbitmq.RabbitMQ
	emailService domainmail.Mail
}

func NewEmailConsumer(
	rabbitMQ *rabbitmq.RabbitMQ,
	emailService domainmail.Mail,
) *EmailConsumer {
	return &EmailConsumer{
		rabbitMQ:     rabbitMQ,
		emailService: emailService,
	}
}

func (consumer *EmailConsumer) Start() error {
	return consumer.rabbitMQ.ConsumeMessages(constants.Events.SendEmail, consumer.handleMessage)
}

func (consumer *EmailConsumer) handleMessage(body []byte) error {
	var msg struct {
		ToEmail      string      `json:"toEmail"`
		Subject      string      `json:"subject"`
		TemplateFile string      `json:"templateFile"`
		Data         interface{} `json:"data"`
	}
	if err := json.Unmarshal(body, &msg); err != nil {
		slog.Error("Failed to unmarshal email notification message", "err", err)
	}

	if err := consumer.emailService.SendEmail(msg.ToEmail, msg.Subject, msg.TemplateFile, msg.Data); err != nil {
		slog.Error("Failed to send email", "err", err)
	}
	return nil
}
