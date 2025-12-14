package consumers

import (
	"encoding/json"
	"hona/backend/bootstrap"
	"hona/backend/internal/infrastructure/mail"
	"hona/backend/internal/infrastructure/rabbitmq"
	"log"
)

var constants = bootstrap.Run().Constants.RabbitMQConstants

type EmailConsumer struct {
	rabbitMQ     *rabbitmq.RabbitMQ
	emailService mail.EmailService
}

func NewEmailConsumer(
	rabbitMQ *rabbitmq.RabbitMQ,
	emailService mail.EmailService,
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
		log.Printf("Failed to unmarshal email notification message: %v", err)
	}

	if err := consumer.emailService.SendEmail(msg.ToEmail, msg.Subject, msg.TemplateFile, msg.Data); err != nil {
		log.Printf("Failed to send email: %v", err)
	}
	return nil
}
