package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"os/exec"
	"time"

	gomail "github.com/Shopify/gomail"
	"github.com/gocarina/gocsv"
)

func SendMember(member PISignUp) {
	var file = WriteCSV(member)
	// Create a new message
	m := gomail.NewMessage()

	// Set sender and recipient
	m.SetHeader("From", "signup@svpromtpusimperii.nl")
	m.SetHeader("To", "secretaris@svpromptusimperii.nl")

	// Set subject and body
	m.SetHeader("Subject", fmt.Sprintf("[Server] Nieuwe aanmelding lid: %s", getFullName(member)))
	m.SetBody("text/plain", "Nieuw lid aangemeld, zie bijlage.")
	m.Attach(file)
	submitMail(m)
}

func WriteCSV(member PISignUp) string {
	// Open the CSV file for writing (folder /app/inschrijvingen is created in the dockerfile)
	var filename = fmt.Sprintf("/app/inschrijvingen/pisignup-%s_%s.csv", member.Surname, time.Now().Format("01-02-2006-15:04:05"))
	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Create a CSV writer
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write CSV header
	// err = gocsv.MarshalCsv(&PISignUp{}, writer)
	// if err != nil {
	// 	panic(err)
	// }
	array := []*PISignUpExport{}
	array = append(array, member.ToPISignUpExport())

	err = gocsv.MarshalFile(array, file)
	if err != nil {
		panic(err)
	}

	println("CSV file created successfully")
	return filename
}

func getFullName(member PISignUp) string {
	var fullName string
	if member.Infix != "" {
		fullName = member.Nickname + " " + member.Infix + " " + member.Surname
	} else {
		fullName = member.Nickname + " " + member.Surname
	}
	return fullName
}

func submitMail(m *gomail.Message) (err error) {
	const sendmail = "/usr/bin/msmtp"
	cmd := exec.Command(sendmail, "-t")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	pw, err := cmd.StdinPipe()
	if err != nil {
		fmt.Println("error", err)
		return
	}

	err = cmd.Start()
	if err != nil {
		fmt.Println("error", err)
		return
	}

	var errs [3]error
	_, errs[0] = m.WriteTo(pw)
	errs[1] = pw.Close()
	errs[2] = cmd.Wait()
	for _, err = range errs {
		if err != nil {
			fmt.Println("error", err)
			return
		}
	}
	fmt.Println("Email sent")
	return
}
