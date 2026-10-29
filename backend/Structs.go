package main

import (
	"regexp"
	"strings"
)

var PlainPhoneNumberRegex *regexp.Regexp = regexp.MustCompile(`^\+[0-9 ]+$`)

type EmailRequest struct {
	Altcha string `json:"altcha"`
}

type ServerEmailCredentials struct {
	email    string
	password string
}

type PISignUp struct {
	LegalFirstNames             string `json:"legal_first_names"`
	Nickname                    string `json:"nickname"`
	Infix                       string `json:"infix"`
	Surname                     string `json:"surname"`
	Phone                       string `json:"phone"`
	DateOfBirth                 string `json:"date_of_birth"`
	Address                     string `json:"address"`
	PostalCode                  string `json:"postal_code"`
	City                        string `json:"city"`
	Country                     string `json:"country"`
	Email                       string `json:"email"`
	Education                   string `json:"education"`
	CohortYear                  string `json:"cohort_year"`
	EmergencyContactFirstName   string `json:"emergency_contact_first_name"`
	EmergencyContactInfix       string `json:"emergency_contact_infix"`
	EmergencyContactSurname     string `json:"emergency_contact_surname"`
	EmergencyContactPhoneNumber string `json:"emergency_contact_phone_number"`
	IBAN                        string `json:"iban"`
	AccountHolder               string `json:"account_holder"`
	Contribution                string `json:"accept_contribution"`
	ApprovalTermsAndConditions  string `json:"accept_terms_and_conditions"`
	Altcha                      string `json:"altcha"`
}

type PISignUpExport struct {
	Voornamen                 string `csv:"Voornamen"`
	Roepnaam                  string `csv:"Roepnaam"`
	Tussenvoegsel             string `csv:"Tussenvoegsel"`
	Achternaam                string `csv:"Achternaam"`
	Geboortedatum             string `csv:"Geboortedatum"`
	Adres                     string `csv:"Adres"`
	Postcode                  string `csv:"Postcode"`
	Woonplaats                string `csv:"Woonplaats"`
	Land                      string `csv:"Land"`
	E_mail                    string `csv:"E-mail"`
	Opleiding                 string `csv:"Opleiding"`
	Type                      string `csv:"Type"`
	Telefoonnummer            string `csv:"Telefoonnummer"`
	Cohortjaar                string `csv:"Cohortjaar"`
	Noodnummer_Naam           string `csv:"Noodnummer (Naam)"`
	Noodnummer_Telefoonnummer string `csv:"Noodnummer (Telefoonnummer)"`
	IBAN                      string `csv:"IBAN"`
	Naam_rekeninghouder       string `csv:"Naam rekeninghouder"`
}

func (member *PISignUp) ToPISignUpExport() *PISignUpExport {
	// Convert TI and I to Technische Informatica and Informatica
	if member.Education == "TI" {
		member.Education = "Technische Informatica"
	} else if member.Education == "I" {
		member.Education = "Informatica"
	}
	// Convert NL and BE to Nederland and België
	if member.Country == "NL" {
		member.Country = "Nederland"
	} else if member.Country == "BE" {
		member.Country = "België"
	}

	var emergencyContactName string
	if member.EmergencyContactInfix != "" {
		emergencyContactName = member.EmergencyContactFirstName + " " + member.EmergencyContactInfix + " " + member.EmergencyContactSurname
	} else {
		emergencyContactName = member.EmergencyContactFirstName + " " + member.EmergencyContactSurname
	}

	// Type is always "Lid"
	return &PISignUpExport{
		Voornamen:                 escapeCSVFormula(member.LegalFirstNames),
		Roepnaam:                  escapeCSVFormula(member.Nickname),
		Tussenvoegsel:             escapeCSVFormula(member.Infix),
		Achternaam:                escapeCSVFormula(member.Surname),
		Geboortedatum:             escapeCSVFormula(member.DateOfBirth),
		Adres:                     escapeCSVFormula(member.Address),
		Postcode:                  escapeCSVFormula(member.PostalCode),
		Woonplaats:                escapeCSVFormula(member.City),
		Land:                      escapeCSVFormula(member.Country),
		E_mail:                    escapeCSVFormula(member.Email),
		Opleiding:                 escapeCSVFormula(member.Education),
		Type:                      "Lid",
		Telefoonnummer:            escapeCSVFormula(member.Phone),
		Cohortjaar:                escapeCSVFormula(member.CohortYear),
		Noodnummer_Naam:           escapeCSVFormula(emergencyContactName),
		Noodnummer_Telefoonnummer: escapeCSVFormula(member.EmergencyContactPhoneNumber),
		IBAN:                      escapeCSVFormula(member.IBAN),
		Naam_rekeninghouder:       escapeCSVFormula(member.AccountHolder),
	}
}

// Excel and LibreOffice execute cells starting with = + - @ as formulas when the CSV is opened.
// Prefixing those with ' makes them plain text, so a signup can't smuggle a formula into the export.
// Phone numbers like +31612345678 are left alone, as they are harmless and are expected to start with +.
func escapeCSVFormula(value string) string {
	if value == "" || PlainPhoneNumberRegex.MatchString(value) {
		return value
	}
	if strings.ContainsRune("=+-@\t\r", rune(value[0])) {
		return "'" + value
	}
	return value
}

type IBANValidationResponse struct {
	Valid        bool              `json:"valid"`
	Messages     []string          `json:"messages"`
	IBAN         string            `json:"iban"`
	BankData     map[string]string `json:"bankData"`
	CheckResults map[string]any    `json:"checkResults"`
}
