package mail

type Mail interface {
	SendEmail(to string, subject string, templateFileName string, data interface{}) error
}