package mail

import (
	"context"
	"fmt"
	"hona/backend/bootstrap"
	"net/url"
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
		mail.WithTLSPortPolicy(mail.TLSMandatory),
	)

	if err != nil {
		panic(err)
	}

	return &EmailService{
		client: client,
	}
}

func (s *EmailService) SendMLEmail(to string, token string) error {
	m := mail.NewMsg()

	if err := m.From(bootstrap.Run().Env.EmailConfig.From); err != nil {
		return err
	}

	if err := m.To(to); err != nil {
		return err
	}

	m.Subject("Email Verification For Pet Yar")

	// Create nice HTML email body
	htmlBody := s.createMagicLinkHTML(token)
	m.SetBodyString(mail.TypeTextHTML, htmlBody)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.client.DialAndSendWithContext(ctx, m); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

func (s *EmailService) createMagicLinkHTML(token string) string {
	appURL := "http://localhost:8080/v1"
	magicLink := fmt.Sprintf("%s/auth/verify?token=%s", appURL, url.QueryEscape(token))

	return fmt.Sprintf(`<!DOCTYPE html>
<html dir="rtl" lang="fa">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <style>
        @import url('https://fonts.googleapis.com/css2?family=Vazirmatn:wght@400;500;700&display=swap');
    </style>
</head>
<body style="font-family: 'Vazirmatn', Arial, sans-serif; margin: 0; padding: 0; background-color: #f5f5f5; direction: rtl;">
    <table role="presentation" style="width: 100%%; border-collapse: collapse;">
        <tr>
            <td align="center" style="padding: 40px 20px;">
                <table role="presentation" style="max-width: 600px; width: 100%%; background-color: #ffffff; border-radius: 20px; overflow: hidden; box-shadow: 0 4px 20px rgba(0,0,0,0.08);">
                    
                    <!-- Header with Orange Gradient -->
                    <tr>
                        <td style="background: linear-gradient(135deg, #FFA726 0%%, #FF8A50 100%%); padding: 50px 40px; text-align: center;">
                            <div style="background-color: white; width: 80px; height: 80px; border-radius: 50%%; margin: 0 auto 20px; display: flex; align-items: center; justify-content: center; box-shadow: 0 4px 15px rgba(0,0,0,0.1);">
                                <span style="font-size: 40px;">🔐</span>
                            </div>
                            <h1 style="color: white; margin: 0; font-size: 28px; font-weight: 700;">
                                ورود به حساب کاربری
                            </h1>
                        </td>
                    </tr>
                    
                    <!-- Content -->
                    <tr>
                        <td style="padding: 40px;">
                            <p style="color: #424242; font-size: 16px; line-height: 28px; margin: 0 0 30px 0; text-align: center;">
                                سلام! 👋
                            </p>
                            <p style="color: #616161; font-size: 15px; line-height: 26px; margin: 0 0 30px 0; text-align: center;">
                                برای ورود به حساب کاربری خود روی دکمه زیر کلیک کنید
                            </p>
                            
                            <!-- Button -->
                            <div style="text-align: center; margin: 35px 0;">
                                <a href="%s" style="background: linear-gradient(135deg, #FFA726 0%%, #FF8A50 100%%); color: white; padding: 16px 50px; text-decoration: none; border-radius: 50px; display: inline-block; font-weight: 700; font-size: 16px; box-shadow: 0 4px 15px rgba(255, 167, 38, 0.3); transition: all 0.3s;">
                                    ورود به حساب
                                </a>
                            </div>
                            
                            <!-- Timer Warning -->
                            <div style="background-color: #FFF3E0; border-right: 4px solid #FFA726; padding: 15px 20px; border-radius: 8px; margin: 30px 0;">
                                <p style="color: #E65100; font-size: 14px; margin: 0; font-weight: 500;">
                                    ⏱ این لینک تا <strong>۱۵ دقیقه</strong> دیگر معتبر است
                                </p>
                            </div>
                            
                            <!-- Alternative Link -->
                            <div style="background-color: #FAFAFA; padding: 20px; border-radius: 12px; margin-top: 30px; border: 1px solid #EEEEEE;">
                                <p style="color: #757575; font-size: 13px; margin: 0 0 10px 0; text-align: center;">
                                    اگر دکمه کار نکرد، لینک زیر را کپی کنید:
                                </p>
                                <p style="color: #FFA726; font-size: 12px; word-break: break-all; margin: 0; text-align: center; direction: ltr;">
                                    %s
                                </p>
                            </div>
                        </td>
                    </tr>
                    
                    <!-- Footer -->
                    <tr>
                        <td style="background-color: #FAFAFA; padding: 30px 40px; text-align: center; border-top: 1px solid #EEEEEE;">
                            <p style="color: #9E9E9E; font-size: 13px; line-height: 22px; margin: 0;">
                                اگر شما درخواست این ایمیل را نداده‌اید، می‌توانید آن را نادیده بگیرید
                            </p>
                            <p style="color: #BDBDBD; font-size: 12px; margin: 15px 0 0 0;">
                                © ۱۴۰۴ تمامی حقوق محفوظ است
                            </p>
                        </td>
                    </tr>
                </table>
            </td>
        </tr>
    </table>
</body>
</html>`, magicLink, magicLink)
}
