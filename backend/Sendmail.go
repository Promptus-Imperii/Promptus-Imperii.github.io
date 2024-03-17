package main

import (
	"fmt"
	"os"
	"os/exec"

	"gopkg.in/gomail.v2"
)

func testSendMail() {
	// Create a new message
	m := gomail.NewMessage()

	// Set sender and recipient
	m.SetHeader("From", "signup@svpromtpusimperii.nl")
	m.SetHeader("To", "secretaris@svpromptusimperii.nl")

	// Set subject and body
	m.SetHeader("Subject", "[Server] Nieuwe aanmelding lid")
	m.SetBody("text/plain", "This is a test email sent using gomail.")
	submitMail(m)
}

func submitMail(m *gomail.Message) (err error) {
	const sendmail = "/usr/bin/sendemail"
	cmd := exec.Command(sendmail, "-t")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	pw, err := cmd.StdinPipe()
	if err != nil {
		fmt.Println("error")
		return
	}

	err = cmd.Start()
	if err != nil {
		fmt.Println("error")
		return
	}

	var errs [3]error
	_, errs[0] = m.WriteTo(pw)
	errs[1] = pw.Close()
	errs[2] = cmd.Wait()
	for _, err = range errs {
		if err != nil {
			fmt.Println("error")
			return
		}
	}
	fmt.Println("send")
	return
}
