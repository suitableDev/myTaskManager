package helper

import (
	"log"
	"os"

	"github.com/keighl/postmark"
)

// SendEmail - Generic function to send emails via Postmark
func SendEmail(to, subject, body string) error {
	client := postmark.NewClient(os.Getenv("POSTMARK_SERVER_TOKEN"), os.Getenv("POSTMARK_ACCOUNT_TOKEN"))
	senderEmail := os.Getenv("POSTMARK_SENDER_EMAIL")

	emailMessage := postmark.Email{
		From:     senderEmail,
		To:       to,
		Subject:  subject,
		TextBody: body,
	}

	_, err := client.SendEmail(emailMessage)
	if err != nil {
		log.Println("Error sending email:", err)
	}
	return err
}

// SendVerificationEmail - Sends verification email
func SendVerificationEmail(email, token string) error {
	body := "Click this link to verify your email: " + os.Getenv("POSTMARK_EMAIL_LINK_ADDRESS") + "/verify?token=" + token
	return SendEmail(email, "Verify Your Email", body)
}

// SendPasswordResetEmail - Sends password reset email
func SendPasswordResetEmail(email, token string) error {
	body := "Click this link to reset your password: " + os.Getenv("POSTMARK_EMAIL_LINK_ADDRESS") + "/users/reset-password?token=" + token
	return SendEmail(email, "Reset Your Password", body)
}
