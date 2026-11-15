package main

import (
	"cmp"
	"crypto/rand"
	"errors"
	"strconv"
)

type Config struct {
	Dev          bool
	DevProxy     string
	SMTPHost     string
	SMTPPort     int
	SMTPFrom     string
	ContactEmail string
	AltchaKey    []byte
}

func loadConfig(getenv func(string) string) (Config, error) {
	c := Config{
		Dev:          getenv("DEV") == "1",
		DevProxy:     cmp.Or(getenv("DEV_PROXY"), "http://localhost:1313"),
		SMTPHost:     getenv("SMTP_HOST"),
		SMTPFrom:     getenv("SERVER_EMAIL_ADDRESS"),
		ContactEmail: getenv("CORRESPONDANCE_EMAIL_ADDRESS"),
		AltchaKey:    []byte(getenv("ALTCHA_HMAC_KEY")),
	}

	port, err := strconv.Atoi(cmp.Or(getenv("SMTP_PORT"), "25"))
	if err != nil {
		return c, errors.New("SMTP_PORT is not a number")
	}
	c.SMTPPort = port

	// In development nothing else needs to be set: `DEV=1 go run .` is enough.
	if c.Dev {
		c.ContactEmail = cmp.Or(c.ContactEmail, "info@example.org")
	}
	if c.ContactEmail == "" {
		return c, errors.New("CORRESPONDANCE_EMAIL_ADDRESS is not set")
	}
	if !c.Dev && (c.SMTPHost == "" || c.SMTPFrom == "") {
		return c, errors.New("SMTP_HOST and SERVER_EMAIL_ADDRESS must be set")
	}

	if !c.Dev && len(c.AltchaKey) == 0 {
		return c, errors.New("ALTCHA_HMAC_KEY is not set")
	}
	// In development, a random key is fine: challenges stop working on restart.
	if len(c.AltchaKey) == 0 {
		c.AltchaKey = make([]byte, 32)
		rand.Read(c.AltchaKey)
	}
	return c, nil
}
