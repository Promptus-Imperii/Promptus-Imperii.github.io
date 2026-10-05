package main

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
	// Convert TI and I to Technische Informatica and Informatica
	if member.Country == "NL" {
		member.Country = "Nederland"
	} else if member.Education == "BE" {
		member.Country = "België"
	}

	var emergencyContactName string
	if member.Infix != "" {
		emergencyContactName = member.EmergencyContactFirstName + " " + member.EmergencyContactInfix + " " + member.EmergencyContactSurname
	} else {
		emergencyContactName = member.EmergencyContactFirstName + " " + member.EmergencyContactSurname
	}

	// Type is always "Lid"
	return &PISignUpExport{
		Voornamen:                 member.LegalFirstNames,
		Roepnaam:                  member.Nickname,
		Tussenvoegsel:             member.Infix,
		Achternaam:                member.Surname,
		Geboortedatum:             member.DateOfBirth,
		Adres:                     member.Address,
		Postcode:                  member.PostalCode,
		Woonplaats:                member.City,
		Land:                      member.Country,
		E_mail:                    member.Email,
		Opleiding:                 member.Education,
		Type:                      "Lid",
		Telefoonnummer:            member.Phone,
		Cohortjaar:                member.CohortYear,
		Noodnummer_Naam:           emergencyContactName,
		Noodnummer_Telefoonnummer: member.EmergencyContactPhoneNumber,
		IBAN:                      member.IBAN,
		Naam_rekeninghouder:       member.AccountHolder,
	}
}

type IBANValidationResponse struct {
	Valid        bool              `json:"valid"`
	Messages     []string          `json:"messages"`
	IBAN         string            `json:"iban"`
	BankData     map[string]string `json:"bankData"`
	CheckResults map[string]any    `json:"checkResults"`
}
