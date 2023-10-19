package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"regexp"
	"strings"

	"github.com/nyaruka/phonenumbers"
)

// ---
// regexes
// ---
var (
	// implemented to https://nl.wikipedia.org/wiki/Postcode#Postcodes_in_Nederland
	//
	// go-staticcheck is angry here, but this is PCRE-compliant
	PostalCodeRegex *regexp.Regexp = regexp.MustCompile(`^[1-9]\d{3}\w{2}$`)
	// this does not need to be fancy, just needs to check if it is somewhat valid.
	//
	// checks for `*@*.*`, which... complies with _some_ RFC at least... right...?
	EmailRegex *regexp.Regexp = regexp.MustCompile(`[\w\d]+@[\w\d]+[.][\w]+`) // jesus christ this is a terrible regex
)

// ---
// validation functions
// ---

// this function takes a postal code and throws a regex at it to see if it is a dutch postal code
//
// it is possible to first check if it is just 4 numbers, making it belgian, but that might be beyond the scope
func validatePostalCode(code string) error {
	if PostalCodeRegex.FindString(code) == "" {
		return errors.New("Postcode is onjuist, probeer het zo: 4818 AJ.")
	}

	if code[len(code)-2] == 'S' && strings.ContainsAny(string(code[len(code)-1]), "ADS") {
		return errors.New("Onjuiste postcode.")
	}

	return nil
}

// TODO make this function support non-dutch phone numbers too
func validatePhoneNumber(numberString string) error {
	// assume the number is dutch in the first place
	number, err := phonenumbers.Parse(numberString, "NL")
	if err != nil {
		return errors.New("Dit is geen telefoonnummer.")
	}

	if !phonenumbers.IsValidNumberForRegion(number, "NL") {
		return errors.New("Geen geldig Nederlands nummer.")
	}

	return nil
}

type IBANValidationResponse struct {
	Valid        bool              `json:"valid"`
	Messages     []string          `json:"messages"`
	IBAN         string            `json:"iban"`
	BankData     map[string]string `json:"bankData"`
	CheckResults map[string]any    `json:"checkResults"`
}

// this function takes an IBAN _without_ spaces
// it then contacts https://openiban.com to check if the IBAN is valid.
// do note that this API currently only supports the following countries
//
// - Belgium
// - Germany
// - Netherlands
// - Luxembourg
// - Switzerland
// - Austria
// - Liechtenstein
func validateIBAN(iban string) error {

	resp, err := http.Get("https://openiban.com/validate/" + iban)
	if err != nil {
		return errors.New("Kon IBAN niet valideren.")
	}
	var ibanval IBANValidationResponse

	err = json.NewDecoder(resp.Body).Decode(&ibanval)
	if err != nil {
		return errors.New("Kon IBAN niet valideren (fout bij externe service).")
	}

	if ibanval.Valid {
		return nil
	}

	return errors.New("IBAN is ongeldig: controleer of je alles goed hebt overgenomen.")
}

// DISCUSS: should this send the email already?
func validateEmail(email string) error {
	if EmailRegex.FindString(email) == "" {
		return errors.New("Email is ongeldig.")
	}

	// idk send the verification email?
	// definitely use a goroutine for that
	//
	// go sendVerificationEmail(email)

	return nil
}
