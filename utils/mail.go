package utils

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/EloToJaa/quizer/db"
	"github.com/EloToJaa/quizer/enum"
	"github.com/EloToJaa/quizer/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/gomail.v2"
)

func ConfirmEmail(userModel *models.User) error {
	from := os.Getenv("EMAIL_ADDRESS")
	password := os.Getenv("EMAIL_PASSWORD")
	smtpServer := os.Getenv("EMAIL_SERVER")
	smtpPort, _ := strconv.Atoi(os.Getenv("EMAIL_PORT"))

	now := time.Now()
	verifyModel := &models.Verify{
		UserId:    userModel.Id,
		CreatedAt: now.Unix(),
		ExpiresAt: now.Add(time.Minute * 15).Unix(),
	}

	verifyCollection := db.GetCollection(enum.Verify)
	result, err := verifyCollection.InsertOne(context.Background(), verifyModel)
	if err != nil {
		return err
	}
	verifyModel.Id = result.InsertedID.(primitive.ObjectID).Hex()

	data := struct {
		Username string
		Url      string
	}{
		Username: userModel.Username,
		Url:      fmt.Sprintf("http://localhost:3000/api/v1/auth/verify/%s", verifyModel.Id),
	}

	tmpl := template.Must(template.ParseFiles("templates/verify.html"))
	var emailBody bytes.Buffer

	if err := tmpl.Execute(&emailBody, data); err != nil {
		log.Fatal(err)
	}

	msg := gomail.NewMessage()
	msg.SetHeader("From", from)
	msg.SetHeader("To", userModel.Email)
	msg.SetHeader("Subject", "Confirm your email")
	msg.SetBody("text/html", emailBody.String())

	dialer := gomail.NewDialer(smtpServer, smtpPort, from, password)
	if err := dialer.DialAndSend(msg); err != nil {
		return err
	}

	return nil
}

func ResetPassword(userModel *models.User) error {
	from := os.Getenv("EMAIL_ADDRESS")
	password := os.Getenv("EMAIL_PASSWORD")
	smtpServer := os.Getenv("EMAIL_SERVER")
	smtpPort, _ := strconv.Atoi(os.Getenv("EMAIL_PORT"))

	now := time.Now()
	resetPasswordModel := &models.ResetPassword{
		UserId:    userModel.Id,
		CreatedAt: now.Unix(),
		ExpiresAt: now.Add(time.Minute * 15).Unix(),
	}

	resetPasswordCollection := db.GetCollection(enum.ResetPassword)
	result, err := resetPasswordCollection.InsertOne(context.Background(), resetPasswordModel)
	if err != nil {
		return err
	}
	resetPasswordModel.Id = result.InsertedID.(primitive.ObjectID).Hex()

	data := struct {
		Username string
		Url      string
	}{
		Username: userModel.Username,
		Url:      fmt.Sprintf("http://localhost:3000/api/v1/auth/reset_password/%s", resetPasswordModel.Id),
	}

	tmpl := template.Must(template.ParseFiles("templates/passwordReset.html"))
	var emailBody bytes.Buffer

	if err := tmpl.Execute(&emailBody, data); err != nil {
		log.Fatal(err)
	}

	msg := gomail.NewMessage()
	msg.SetHeader("From", from)
	msg.SetHeader("To", userModel.Email)
	msg.SetHeader("Subject", "Reset your password")
	msg.SetBody("text/html", emailBody.String())

	dialer := gomail.NewDialer(smtpServer, smtpPort, from, password)
	if err := dialer.DialAndSend(msg); err != nil {
		return err
	}

	return nil
}
