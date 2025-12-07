package utils

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"

	"gopkg.in/gomail.v2"
)

func GenerateOTP() string {
	rand.Seed(time.Now().UnixNano())
	otp := rand.Intn(900000) + 100000 // Generate 6-digit OTP
	return strconv.Itoa(otp)
}

func SendOTPEmail(toEmail, otp string) error {
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))
	smtpEmail := os.Getenv("SMTP_EMAIL")
	smtpPassword := os.Getenv("SMTP_PASSWORD")

	if smtpEmail == "" || smtpPassword == "" {
		// In development, just log the OTP
		fmt.Printf("📧 OTP for %s: %s\n", toEmail, otp)
		return nil
	}

	m := gomail.NewMessage()
	m.SetHeader("From", smtpEmail)
	m.SetHeader("To", toEmail)
	m.SetHeader("Subject", "Verify Your Email - Crypto Wallet")
	
	body := fmt.Sprintf(`
		<html>
		<body>
			<h2>Welcome to Crypto Wallet!</h2>
			<p>Your OTP for email verification is:</p>
			<h1 style="color: #4F46E5; font-size: 32px;">%s</h1>
			<p>This OTP will expire in 10 minutes.</p>
			<p>If you didn't request this, please ignore this email.</p>
		</body>
		</html>
	`, otp)
	
	m.SetBody("text/html", body)

	d := gomail.NewDialer(smtpHost, smtpPort, smtpEmail, smtpPassword)

	if err := d.DialAndSend(m); err != nil {
		// Log but don't fail in development
		fmt.Printf("⚠️  Email send failed (using console instead): %v\n", err)
		fmt.Printf("📧 OTP for %s: %s\n", toEmail, otp)
		return nil
	}

	return nil
}
