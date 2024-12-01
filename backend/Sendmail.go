package main

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"log"

	gomail "github.com/Shopify/gomail"
	"github.com/gin-gonic/gin"
	"github.com/gocarina/gocsv"
)

func SendMemberInfoEmail(member PISignUp, serverEmailCredentials ServerEmailCredentials, correspondanceEmail string) error {
	if gin.Mode() == gin.TestMode {
		log.Println("Testing mode: email will not be sent")
		return nil
	}

	// Write member info to a CSV file
	csv_bytes, err := WriteToCSV(member)
	if err != nil {
		return err
	}

	// Create a new message
	m := gomail.NewMessage()

	// Set sender and recipient
	m.SetHeader("From", serverEmailCredentials.email)
	m.SetHeader("To", correspondanceEmail)

	// Set subject and body
	m.SetHeader("Subject", fmt.Sprintf("[Server] Nieuwe aanmelding lid: %s", getFullName(member)))
	m.SetBody("text/plain", "Nieuw lid aangemeld, zie bijlage.")
	m.AttachReader("nieuw_lid.csv", bytes.NewReader(csv_bytes))

	err = SendEmail(serverEmailCredentials, m)

	if err != nil {
		log.Println("Error sending email to contact email ", err)
		return err
	}

	log.Println("Email to contact sent")
	return nil
}

func SendNotificationEmail(member PISignUp, serverEmailCredentials ServerEmailCredentials, correspondanceEmail string) error {
	if gin.Mode() == gin.TestMode || gin.Mode() == gin.DebugMode {
		log.Println("Testing or debug mode: email will not be sent")
		return nil
	}

	// Create a new message
	m := gomail.NewMessage()

	// Set sender and recipient
	m.SetHeader("From", serverEmailCredentials.email)
	m.SetHeader("To", member.Email)

	// Set subject and body
	m.SetHeader("Subject", "No-reply: Bevestiging aanmelding S.V Promptus Imperii.")
	m.SetBody("text/plain", fmt.Sprintf("Beste,\n"+
		"\n"+
		"Bedankt voor je aanmelding bij S.V Promptus Imperii. De secretaris zal jouw aanmelding zo snel mogelijk in behandeling nemen. Dit kan een paar dagen duren, aangezien het een handmatig proces is.\n"+
		"Als je na een week nog steeds niets gehoord hebt, aarzel dan niet om contact op te nemen met %s.", correspondanceEmail))

	err := SendEmail(serverEmailCredentials, m)

	if err != nil {
		log.Println("Error writing confirmation email to ", member.Email, err)
		return err
	}
	log.Println("Confirmation email sent")
	return nil
}

func SendEmail(serverEmailCredentials ServerEmailCredentials, message *gomail.Message) error {

	// Send via server email/office365
	d := gomail.NewDialer("smtp.office365.com", 587, serverEmailCredentials.email, serverEmailCredentials.password)

	d.TLSConfig = &tls.Config{ServerName: "smtp.office365.com"}
	err := d.DialAndSend(message)

	if err != nil {
		return err
	}

	return nil
}

func WriteToCSV(member PISignUp) ([]byte, error) {
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

	if member.Nickname != "" {
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
