package main

import "time"

type PISignUP struct {
	LegalFirstNames  string    `json:"legalfirstnames"`
	Member           Contact   `json:"member"`
	DateOfBirth      time.Time `json:"date_of_birth"`
	Address          string    `json:"address"`
	PostalCode       string    `json:"postal_code"`
	City             string    `json:"city"`
	Email            string    `json:"email"`
	Course           string    `json:"course"`
	Cohort           string    `json:"cohort"`
	EmergencyContact Contact   `json:"emergency_contact"`
	IBAN             string    `json:"iban"`
	AccountHolder    string    `json:"account_holder"`
}

type Contact struct {
	FirstName   string `json:"firstname"`
	Infix       string `json:"infix"` // tussenvoegsel (de, van, den etc.)
	LastName    string `json:"lastname"`
	PhoneNumber string `json:"phone"`
}

type IBANValidationResponse struct {
	Valid        bool              `json:"valid"`
	Messages     []string          `json:"messages"`
	IBAN         string            `json:"iban"`
	BankData     map[string]string `json:"bankData"`
	CheckResults map[string]any    `json:"checkResults"`
}
