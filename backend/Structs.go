package main

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

type IBANValidationResponse struct {
	Valid        bool              `json:"valid"`
	Messages     []string          `json:"messages"`
	IBAN         string            `json:"iban"`
	BankData     map[string]string `json:"bankData"`
	CheckResults map[string]any    `json:"checkResults"`
}
