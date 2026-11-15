package main

import (
	"net/url"
	"strings"
	"testing"
)

func validForm() url.Values {
	return url.Values{
		"legal_first_names":              {"Zoë Johanna"},
		"nickname":                       {"Zoë"},
		"infix":                          {"van"},
		"surname":                        {"Müller"},
		"date_of_birth":                  {"2000-10-12"},
		"address":                        {"Lovensdijkstraat 61"},
		"postal_code":                    {"4818aj"},
		"city":                           {"Breda"},
		"country":                        {"NL"},
		"email":                          {"zoe@example.org"},
		"phone":                          {"+31 6-1234 5678"},
		"cohort_year":                    {"2025/2026"},
		"emergency_contact_first_name":   {"Greetje"},
		"emergency_contact_infix":        {"de"},
		"emergency_contact_surname":      {"Vries"},
		"emergency_contact_phone_number": {"+31687654321"},
		"iban":                           {"nl91 abna 0417 1643 00"},
		"account_holder":                 {"Z. J. van Müller"},
		"accept_contribution":            {"on"},
		"accept_terms_and_conditions":    {"on"},
	}
}

func TestParseSignupValidNormalizes(t *testing.T) {
	s, errs := parseSignup(validForm(), testNow)
	if len(errs) != 0 {
		t.Fatalf("errors: %v", errs)
	}
	if s.PostalCode != "4818 AJ" || s.Phone != "+31612345678" || s.IBAN != "NL91ABNA0417164300" {
		t.Errorf("not normalized: %+v", s)
	}
	if s.LegalFirstNames != "Zoë Johanna" || s.EmergencyInfix != "de" {
		t.Errorf("fields: %+v", s)
	}
}

func TestParseSignupBelgium(t *testing.T) {
	f := validForm()
	f.Set("country", "BE")
	f.Set("postal_code", "2000")
	f.Set("iban", "BE68 5390 0754 7034")
	s, errs := parseSignup(f, testNow)
	if len(errs) != 0 || s.PostalCode != "2000" {
		t.Fatalf("errs %v, postal %q", errs, s.PostalCode)
	}
}

func TestParseSignupErrors(t *testing.T) {
	cases := []struct {
		field, value, want string
	}{
		{"legal_first_names", "  ", "Voornamen is verplicht"},
		{"surname", strings.Repeat("a", 201), "Achternaam is te lang"},
		{"city", "Bre\r\nda", "Woonplaats bevat ongeldige tekens"},
		{"country", "DE", "Land is ongeldig"},
		{"postal_code", "1234", "Nederlandse postcode"},
		{"postal_code", "0123AB", "Nederlandse postcode"},
		{"date_of_birth", "2030-01-01", "Geboortedatum is ongeldig"},
		{"date_of_birth", "1899-12-31", "Geboortedatum is ongeldig"},
		{"date_of_birth", "12-10-2000", "Geboortedatum is ongeldig"},
		{"email", "Jan <jan@example.org>", "E-mailadres is ongeldig"},
		{"email", "jan", "E-mailadres is ongeldig"},
		{"phone", "+3161234", "Telefoonnummer is ongeldig"},
		{"phone", "0612345678", "Telefoonnummer is ongeldig"},
		{"phone", "0031612345678", "Telefoonnummer is ongeldig"},
		{"phone", "+31 6 1234 567a", "Telefoonnummer is ongeldig"},
		{"phone", "+31 6 1234 56789", "Telefoonnummer is ongeldig"},
		{"phone", "+32 4 1234", "Telefoonnummer is ongeldig"},
		{"phone", " - ", "Telefoonnummer is verplicht"},
		{"emergency_contact_phone_number", "()", "Telefoonnummer noodcontact is verplicht"},
		{"emergency_contact_phone_number", "+31", "Telefoonnummer noodcontact is ongeldig"},
		{"cohort_year", "2025/2027", "Cohortjaar is ongeldig"},
		{"cohort_year", "25/26", "Cohortjaar is ongeldig"},
		{"iban", "NL91ABNA0417164301", "IBAN is ongeldig"},
		{"iban", "NL91", "IBAN is ongeldig"},
		{"accept_contribution", "", "contributie"},
		{"accept_terms_and_conditions", "", "algemene voorwaarden"},
		{"city", "=HYPERLINK(\"http://evil\")", "Woonplaats mag niet beginnen met"},
		{"surname", "@SUM(A1)", "Achternaam mag niet beginnen met"},
		{"address", "-2+3", "Adres mag niet beginnen met"},
		{"account_holder", "+cmd", "Naam rekeninghouder mag niet beginnen met"},
		{"nickname", " =1", "Roepnaam mag niet beginnen met"},
		{"email", "=x@example.org", "E-mailadres mag niet beginnen met"},
	}
	for _, c := range cases {
		f := validForm()
		f.Set(c.field, c.value)
		_, errs := parseSignup(f, testNow)
		if len(errs) != 1 || !strings.Contains(errs[0], c.want) {
			t.Errorf("%s=%q: got %v, want one error containing %q", c.field, c.value, errs, c.want)
		}
	}
}

func TestParseSignupBelgianPostalCode(t *testing.T) {
	f := validForm()
	f.Set("country", "BE")
	f.Set("postal_code", "20000")
	_, errs := parseSignup(f, testNow)
	if len(errs) != 1 || !strings.Contains(errs[0], "Belgische postcode") {
		t.Errorf("got %v", errs)
	}
}

func TestParseSignupEmptyFormListsAllRequired(t *testing.T) {
	_, errs := parseSignup(url.Values{}, testNow)
	if len(errs) != 17 {
		t.Errorf("got %d errors: %v", len(errs), errs)
	}
}

func TestValidIBAN(t *testing.T) {
	for _, s := range []string{"NL91ABNA0417164300", "BE68539007547034", "DE89370400440532013000"} {
		if !validIBAN(s) {
			t.Errorf("%s rejected", s)
		}
	}
	for _, s := range []string{"", "NL91ABNA0417164301", "XX00", "NL91ABNA04171643!!"} {
		if validIBAN(s) {
			t.Errorf("%s accepted", s)
		}
	}
}

func TestParseSignupPhoneNumbers(t *testing.T) {
	for input, want := range map[string]string{
		"+31 6-1234 5678":   "+31612345678",
		"+31 (0)6 12345678": "+31612345678",
		"+31 10 123 4567":   "+31101234567",
		"+32 467 30 05 12":  "+32467300512",
		"+32 2 123 45 67":   "+3221234567",
		"+49 1512 3456789":  "+4915123456789",
	} {
		f := validForm()
		f.Set("phone", input)
		s, errs := parseSignup(f, testNow)
		if len(errs) != 0 || s.Phone != want {
			t.Errorf("%q: got %q, errors %v", input, s.Phone, errs)
		}
	}
}
