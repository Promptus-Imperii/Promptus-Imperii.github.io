package main

import (
	"bytes"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/gocarina/gocsv"
	"github.com/wneessen/go-mail"
)

func SendMemberInfoEmail(member PISignUp, serverEmailCredentials ServerEmailCredentials, correspondanceEmail string) error {
	if gin.Mode() == gin.TestMode {
		log.Println("Testing mode: email will not be sent")
		return nil
	}

	// Write member info to a CSV file
	csvBytes, err := WriteToCSV(member)
	if err != nil {
		return err
	}

	// Create a new email message
	m := mail.NewMsg()

	// Set sender and recipient
	m.From(serverEmailCredentials.email)
	m.To(correspondanceEmail)

	// Set subject and body
	m.Subject(fmt.Sprintf("[Server] Nieuwe aanmelding lid: %s", getFullName(member)))
	m.SetBodyString(mail.TypeTextPlain, "Nieuw lid aangemeld, zie bijlage.")

	// Attach the CSV file
	m.AttachReader("nieuw_lid.csv", bytes.NewReader(csvBytes))

	// Send the email
	err = SendEmail(serverEmailCredentials, m)
	if err != nil {
		log.Println("Error sending email to contact email:", err)
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

	// Create a new email message
	m := mail.NewMsg()

	// Set sender and recipient
	m.From(serverEmailCredentials.email)
	m.To(member.Email)

	// Set subject and body
	m.Subject("No-reply: Bevestiging aanmelding S.V Promptus Imperii.")
	body := fmt.Sprintf(
		"Beste,\n\nBedankt voor je aanmelding bij S.V Promptus Imperii. De secretaris zal jouw aanmelding zo snel mogelijk in behandeling nemen. Dit kan een paar dagen duren, aangezien het een handmatig proces is.\nAls je na een week nog steeds niets gehoord hebt, aarzel dan niet om contact op te nemen met %s.",
		correspondanceEmail,
	)
	m.SetBodyString(mail.TypeTextPlain, body)

	// Send the email
	err := SendEmail(serverEmailCredentials, m)
	if err != nil {
		log.Println("Error writing confirmation email to", member.Email, err)
		return err
	}

	log.Println("Confirmation email sent")
	return nil
}

func SendEmail(serverEmailCredentials ServerEmailCredentials, message *mail.Msg) error {
	// Configure the email client
	client, err := mail.NewClient(
		"smtp.office365.com",
		mail.WithPort(587),
		mail.WithSMTPAuth(mail.SMTPAuthLogin),
		mail.WithUsername(serverEmailCredentials.email),
		mail.WithPassword(serverEmailCredentials.password),
		mail.WithTLSPortPolicy(mail.TLSMandatory),
	)
	if err != nil {
		return fmt.Errorf("error creating mail client: %w", err)
	}

	// Send the email
	if err := client.DialAndSend(message); err != nil {
		return fmt.Errorf("error sending email: %w", err)
	}

	return nil
}

func WriteToCSV(member PISignUp) ([]byte, error) {
	array := []*PISignUpExport{}
	array = append(array, member.ToPISignUpExport())
	csvBytes, err := gocsv.MarshalBytes(array)

	if err != nil {
		log.Println("Error writing csv:", err)
		return nil, err
	}
	log.Println("CSV file created successfully")
	return csvBytes, nil
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
