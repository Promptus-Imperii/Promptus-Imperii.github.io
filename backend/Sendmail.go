package main

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"log"
	"os"

	gomail "github.com/Shopify/gomail"
	"github.com/gin-gonic/gin"
	"github.com/gocarina/gocsv"
)

func SendMember(member PISignUp) error {
	if gin.Mode() == gin.TestMode {
		log.Println("Testing mode: email will not be sent")
		return nil
	}
	csv_bytes, err := WriteCSV(member)
	if err != nil {
		return err
	}

	// Create a new message
	m := gomail.NewMessage()

	// Set sender and recipient
	m.SetHeader("From", "signup@svpromptusimperii.nl")
	m.SetHeader("To", "secretaris@svpromptusimperii.nl")

	// Set subject and body
	m.SetHeader("Subject", fmt.Sprintf("[Server] Nieuwe aanmelding lid: %s", getFullName(member)))
	m.SetBody("text/plain", "Nieuw lid aangemeld, zie bijlage.")
	m.AttachReader("nieuw_lid.csv", bytes.NewReader(csv_bytes))
	email_password := os.Getenv("EMAIL_PASSWORD")
	d := gomail.NewDialer("smtp.office365.com", 587, "signup@svpromptusimperii.nl", email_password)

	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	err = d.DialAndSend(m)

	if err != nil {
		log.Println("Error writing sending email to", member.Email, err)
		return err
	}
	log.Println("Email sent")
	return nil
}

func WriteCSV(member PISignUp) ([]byte, error) {
	array := []*PISignUpExport{}
	array = append(array, member.ToPISignUpExport())
	csv_bytes, err := gocsv.MarshalBytes(array)

	if err != nil {
		log.Println("Error writing csv:", err)
		return nil, err
	}
	log.Println("CSV file created successfully")
	return csv_bytes, nil
}

func getFullName(member PISignUp) string {
	var fullName string
	var firstName string

	if member.Nickname == "" {
		firstName = member.Nickname
	} else {
		firstName = member.LegalFirstNames
	}

	if member.Infix != "" {
		fullName = firstName + " " + member.Infix + " " + member.Surname
	} else {
		fullName = firstName + " " + member.Surname
	}

	return fullName
}
