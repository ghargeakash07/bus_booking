package utils

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
	"os"

	"akash-project-go/models"

	"github.com/sirupsen/logrus"
)

func BookingBusMailTrigger(userDetail models.UserDetail, bookingDetails models.BookingBus) error {
	smtpHost := os.Getenv("MAIL_HOST")
	smtpPort := os.Getenv("MAIL_PORT")

	smtpUsername := os.Getenv("MAIL_USERNAME")
	smtpPassword := os.Getenv("MAIL_PASSWORD")
	from := os.Getenv("MAIL_FROM_ADDRESS")
	fromName := os.Getenv("MAIL_FROM_NAME")

	// Create the email message
	subject := "Booking Confirmation"

	// Customize the body to exclude created and updated dates
	body := fmt.Sprintf("Dear %s,\n\nThank you for your booking!\n\nBooking Details:\n Booking-Id:%d\n- Name: %s\n- Gender: %s\n- Type: %s\n- Date: %s\n- From: %s\n- To: %s\n- Email: %s\n\nBest Regards,\n%s",
		userDetail.Name,
		bookingDetails.BookingId,
		bookingDetails.PassengerName,
		bookingDetails.Gender,
		bookingDetails.BusType,
		bookingDetails.BookingDate,
		bookingDetails.Location,
		bookingDetails.Destination,
		userDetail.Email,
		fromName)

	message := []byte(fmt.Sprintf("Subject: %s\r\nFrom: %s\r\nTo: %s\r\n\r\n%s", subject, from, userDetail.Email, body))

	// Connect to the SMTP server
	addr := fmt.Sprintf("%s:%s", smtpHost, smtpPort)
	c, err := smtp.Dial(addr)
	if err != nil {
		logrus.Error("Failed to connect to the SMTP server: ", err)
		return err
	}
	defer c.Close()

	// Use TLS if required
	if smtpPort == "587" {
		tlsConfig := &tls.Config{
			InsecureSkipVerify: true,
			ServerName:         smtpHost,
		}
		if err = c.StartTLS(tlsConfig); err != nil {
			logrus.Error("Failed to start TLS: ", err)
			return err
		}
	}

	// Set up authentication
	auth := smtp.PlainAuth("", smtpUsername, smtpPassword, smtpHost)
	if err = c.Auth(auth); err != nil {
		logrus.Error("Failed to authenticate: ", err)
		return err
	}

	// Set the sender and recipient
	if err = c.Mail(from); err != nil {
		logrus.Error("Failed to set the sender address: ", err)
		return err
	}
	if err = c.Rcpt(userDetail.Email); err != nil {
		logrus.Error("Failed to set the recipient address: ", err)
		return err
	}

	// Send the email body
	w, err := c.Data()
	if err != nil {
		logrus.Error("Failed to get the data writer: ", err)
		return err
	}
	defer w.Close()

	_, err = w.Write(message)
	if err != nil {
		logrus.Error("Failed to write the message: ", err)
		return err
	}
	logrus.Info("Email sent successfully to ", userDetail.Email)
	return nil
}

func SendOTPEmail(userDetail models.UserDetail) error {
	smtpHost := os.Getenv("MAIL_HOST")
	smtpPort := os.Getenv("MAIL_PORT")
	smtpUsername := os.Getenv("MAIL_USERNAME")
	smtpPassword := os.Getenv("MAIL_PASSWORD")
	from := os.Getenv("MAIL_FROM_ADDRESS")
	fromName := os.Getenv("MAIL_FROM_NAME")

	// Create the email message
	subject := "Your OTP Code"
	body := fmt.Sprintf("Dear %s,\n\nYour OTP code is: %s\n\nPlease use this code to reset your password.\n\nBest Regards,\n%s",
		userDetail.Name,
		userDetail.OTP, // Assuming OTP is stored in userDetail
		fromName)

	message := []byte(fmt.Sprintf("Subject: %s\r\nFrom: %s\r\nTo: %s\r\n\r\n%s", subject, from, userDetail.Email, body))

	// Connect to the SMTP server
	addr := fmt.Sprintf("%s:%s", smtpHost, smtpPort)
	c, err := smtp.Dial(addr)
	if err != nil {
		logrus.Error("Failed to connect to the SMTP server: ", err)
		return err
	}
	defer c.Close()

	// Use TLS if required
	if smtpPort == "587" {
		tlsConfig := &tls.Config{
			InsecureSkipVerify: true,
			ServerName:         smtpHost,
		}
		if err = c.StartTLS(tlsConfig); err != nil {
			logrus.Error("Failed to start TLS: ", err)
			return err
		}
	}

	// Set up authentication
	auth := smtp.PlainAuth("", smtpUsername, smtpPassword, smtpHost)
	if err = c.Auth(auth); err != nil {
		logrus.Error("Failed to authenticate: ", err)
		return err
	}

	// Set the sender and recipient
	if err = c.Mail(from); err != nil {
		logrus.Error("Failed to set the sender address: ", err)
		return err
	}
	if err = c.Rcpt(userDetail.Email); err != nil {
		logrus.Error("Failed to set the recipient address: ", err)
		return err
	}

	// Send the email body
	w, err := c.Data()
	if err != nil {
		logrus.Error("Failed to get the data writer: ", err)
		return err
	}
	defer w.Close()

	_, err = w.Write(message)
	if err != nil {
		logrus.Error("Failed to write the message: ", err)
		return err
	}

	logrus.Info("OTP email sent successfully to ", userDetail.Email)
	return nil
}
