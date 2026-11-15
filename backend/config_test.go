package main

import "testing"

func env(m map[string]string) func(string) string {
	return func(k string) string { return m[k] }
}

func TestLoadConfigDefaults(t *testing.T) {
	c, err := loadConfig(env(map[string]string{
		"SMTP_HOST":                    "example-org.mail.protection.outlook.com",
		"SERVER_EMAIL_ADDRESS":         "noreply@example.org",
		"CORRESPONDANCE_EMAIL_ADDRESS": "info@example.org",
		"ALTCHA_HMAC_KEY":              "k",
	}))
	if err != nil {
		t.Fatal(err)
	}
	if c.SMTPHost != "example-org.mail.protection.outlook.com" || c.SMTPPort != 25 || c.SMTPFrom != "noreply@example.org" {
		t.Errorf("smtp = %s:%d from %s", c.SMTPHost, c.SMTPPort, c.SMTPFrom)
	}
	if c.Dev || c.DevProxy != "http://localhost:1313" {
		t.Errorf("dev = %v %s", c.Dev, c.DevProxy)
	}
	if string(c.AltchaKey) != "k" {
		t.Errorf("altcha key = %q", c.AltchaKey)
	}
}

func TestLoadConfigErrors(t *testing.T) {
	cases := map[string]map[string]string{
		"no correspondence":  {"SMTP_HOST": "mx.example.org", "SERVER_EMAIL_ADDRESS": "noreply@example.org", "ALTCHA_HMAC_KEY": "k"},
		"no sender":          {"CORRESPONDANCE_EMAIL_ADDRESS": "info@example.org"},
		"relay without host": {"CORRESPONDANCE_EMAIL_ADDRESS": "info@example.org", "SERVER_EMAIL_ADDRESS": "noreply@example.org"},
		"no altcha key":      {"CORRESPONDANCE_EMAIL_ADDRESS": "info@example.org", "SERVER_EMAIL_ADDRESS": "noreply@example.org", "SMTP_HOST": "mx.example.org"},
		"bad port":           {"DEV": "1", "CORRESPONDANCE_EMAIL_ADDRESS": "info@example.org", "SMTP_PORT": "abc"},
	}
	for name, m := range cases {
		if _, err := loadConfig(env(m)); err == nil {
			t.Errorf("%s: expected error", name)
		}
	}
}

func TestLoadConfigDevNeedsNoCredentials(t *testing.T) {
	c, err := loadConfig(env(map[string]string{
		"DEV": "1",
	}))
	if err != nil {
		t.Fatal(err)
	}
	if !c.Dev || len(c.AltchaKey) != 32 || c.ContactEmail != "info@example.org" {
		t.Errorf("got %+v", c)
	}
}
