package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"slices"
	"strings"
	"time"
	"unicode/utf8"
)

type Signup struct {
	LegalFirstNames    string
	Nickname           string
	Infix              string
	Surname            string
	DateOfBirth        string
	Address            string
	PostalCode         string
	City               string
	Country            string
	Email              string
	Phone              string
	CohortYear         string
	EmergencyFirstName string
	EmergencyInfix     string
	EmergencySurname   string
	EmergencyPhone     string
	IBAN               string
	AccountHolder      string
}

func signupRules(now time.Time) []rule {
	return []rule{
		{name: "legal_first_names", label: "Voornamen"},
		{name: "nickname", label: "Roepnaam", optional: true},
		{name: "infix", label: "Tussenvoegsel", optional: true},
		{name: "surname", label: "Achternaam"},
		{name: "date_of_birth", label: "Geboortedatum", check: checkBirthDate(now)},
		{name: "address", label: "Adres"},
		{name: "country", label: "Land", check: checkOneOf("Land", "NL", "BE")},
		{name: "postal_code", label: "Postcode", normalize: upperNoSpaces, check: checkPostalCode},
		{name: "city", label: "Woonplaats"},
		{name: "email", label: "E-mailadres", check: checkEmail},
		{name: "phone", label: "Telefoonnummer", skipFormulaCheck: true, normalize: stripPhoneSeparators, check: checkPhone("Telefoonnummer")},
		{name: "cohort_year", label: "Cohortjaar", check: checkCohortYear},
		{name: "emergency_contact_first_name", label: "Voornaam noodcontact"},
		{name: "emergency_contact_infix", label: "Tussenvoegsel noodcontact", optional: true},
		{name: "emergency_contact_surname", label: "Achternaam noodcontact"},
		{name: "emergency_contact_phone_number", label: "Telefoonnummer noodcontact", skipFormulaCheck: true, normalize: stripPhoneSeparators, check: checkPhone("Telefoonnummer noodcontact")},
		{name: "iban", label: "IBAN", normalize: upperNoSpaces, check: checkIBAN},
		{name: "account_holder", label: "Naam rekeninghouder"},
	}
}

func parseSignup(form url.Values, now time.Time) (Signup, []string) {
	values, errs := validate(form, signupRules(now))
	if form.Get("accept_contribution") == "" {
		errs = append(errs, "Je moet akkoord gaan met de contributie")
	}
	if form.Get("accept_terms_and_conditions") == "" {
		errs = append(errs, "Je moet akkoord gaan met de algemene voorwaarden")
	}
	return Signup{
		LegalFirstNames:    values["legal_first_names"],
		Nickname:           values["nickname"],
		Infix:              values["infix"],
		Surname:            values["surname"],
		DateOfBirth:        values["date_of_birth"],
		Address:            values["address"],
		PostalCode:         values["postal_code"],
		City:               values["city"],
		Country:            values["country"],
		Email:              values["email"],
		Phone:              values["phone"],
		CohortYear:         values["cohort_year"],
		EmergencyFirstName: values["emergency_contact_first_name"],
		EmergencyInfix:     values["emergency_contact_infix"],
		EmergencySurname:   values["emergency_contact_surname"],
		EmergencyPhone:     values["emergency_contact_phone_number"],
		IBAN:               values["iban"],
		AccountHolder:      values["account_holder"],
	}, errs
}

var csvHeader = []string{
	"Voornamen", "Roepnaam", "Tussenvoegsel", "Achternaam", "Geboortedatum", "Adres", "Postcode",
	"Woonplaats", "Land", "E-mail", "Opleiding", "Type", "Telefoonnummer", "Cohortjaar",
	"Noodnummer (Naam)", "Noodnummer (Telefoonnummer)", "IBAN", "Naam rekeninghouder",
}

var countryNames = map[string]string{"NL": "Nederland", "BE": "België"}

// csv builds the secretary's CSV. The BOM makes Excel read UTF-8.
// "Opleiding" stays empty for import compatibility.
func (s Signup) csv() ([]byte, error) {
	var b bytes.Buffer
	b.WriteString("\uFEFF")
	w := csv.NewWriter(&b)
	_ = w.Write(csvHeader)
	_ = w.Write([]string{
		s.LegalFirstNames, s.Nickname, s.Infix, s.Surname, s.DateOfBirth, s.Address, s.PostalCode,
		s.City, countryNames[s.Country], s.Email, "", "Lid", s.Phone, s.CohortYear,
		s.emergencyName(), s.EmergencyPhone, s.IBAN, s.AccountHolder,
	})
	w.Flush() // the writer keeps the first error for w.Error
	return b.Bytes(), w.Error()
}

func (s Signup) fullName() string {
	first := s.Nickname
	if first == "" {
		first = s.LegalFirstNames
	}
	return joinName(first, s.Infix, s.Surname)
}

func (s Signup) emergencyName() string {
	return joinName(s.EmergencyFirstName, s.EmergencyInfix, s.EmergencySurname)
}

func joinName(parts ...string) string {
	return strings.Join(slices.DeleteFunc(parts, func(p string) bool { return p == "" }), " ")
}

const confirmationBody = "Beste,\n\n" +
	"Bedankt voor je aanmelding bij S.V Promptus Imperii. De secretaris zal jouw aanmelding zo snel mogelijk in behandeling nemen. " +
	"Dit kan een paar dagen duren, aangezien het een handmatig proces is.\n" +
	"Als je na een week nog steeds niets gehoord hebt, aarzel dan niet om contact op te nemen met %s."

func (a *app) handleSignup(w http.ResponseWriter, r *http.Request) {
	captcha, ok := verifyAltcha(a.cfg.AltchaKey, r.PostForm.Get("altcha"))
	if !ok {
		writeErrors(w, http.StatusUnprocessableEntity, []string{msgCaptcha})
		return
	}
	member, errs := parseSignup(r.PostForm, a.now())
	if len(errs) > 0 {
		log.Printf("signup rejected: %d validation errors", len(errs)) //nolint:gosec // only logs a count
		writeErrors(w, http.StatusUnprocessableEntity, errs)
		return
	}
	// One signup per solved captcha, so a captcha can't be replayed to mail many addresses.
	if !a.captchas.claim(captcha, a.now()) {
		log.Println("signup rejected: captcha already used")
		writeErrors(w, http.StatusUnprocessableEntity, []string{msgCaptcha})
		return
	}

	failedMsg := []string{"Er is iets fout gegaan tijdens het verwerken van je aanmelding. Meld jezelf aan via " + a.cfg.ContactEmail}
	data, err := member.csv()
	if err != nil {
		log.Printf("signup csv: %v", err)
		a.captchas.release(captcha)
		writeErrors(w, http.StatusInternalServerError, failedMsg)
		return
	}
	err = a.mailer.Send(r.Context(), Mail{
		To:         a.cfg.ContactEmail,
		Subject:    "[Server] Nieuwe aanmelding lid: " + truncate(member.fullName(), maxSubjectName),
		Body:       "Nieuw lid aangemeld, zie bijlage.",
		Attachment: &Attachment{Name: "nieuw_lid.csv", ContentType: "text/csv; charset=utf-8", Data: data},
	})
	if err != nil {
		log.Printf("signup mail to secretary: %v", err)
		a.captchas.release(captcha)
		writeErrors(w, http.StatusInternalServerError, failedMsg)
		return
	}

	// The signup is in, so a failed confirmation is only logged.
	err = a.mailer.Send(r.Context(), Mail{
		To:      member.Email,
		Subject: "No-reply: Bevestiging aanmelding S.V Promptus Imperii.",
		Body:    fmt.Sprintf(confirmationBody, a.cfg.ContactEmail),
	})
	if err != nil {
		log.Printf("signup confirmation mail: %v", err)
	}
	log.Println("signup processed")
	redirect(w, "/signup_success/")
}

// maxSubjectName keeps the subject under the 998-character line limit.
const maxSubjectName = 60

func truncate(s string, n int) string {
	if utf8.RuneCountInString(s) <= n {
		return s
	}
	return string([]rune(s)[:n]) + "…"
}
