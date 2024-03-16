package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/mail"
	"regexp"
	"strings"
	"time"

	"github.com/nyaruka/phonenumbers"
)

// ---
// regexes
// ---
var (
	// Based on https://stackoverflow.com/a/17898538
	// Lookahead is not currently supported in standard library regex, so banned letters are not checked.
	DutchPostalCodeRegex   *regexp.Regexp = regexp.MustCompile(`^[1-9][0-9]{3}[A-Z]{2}$`)
	BelgianPostalCodeRegex *regexp.Regexp = regexp.MustCompile(`^\d{4}$`)
)

// ---
// validation functions
// ---

// this function takes a postal code and throws a regex at it to see if it is a dutch postal code
//
// it is possible to first check if it is just 4 numbers, making it belgian, but that might be beyond the scope
func validatePostalCode(postalCode string) error {
	// Normalize postal codes: capitalize all letters and remove all spaces
	postalCode = strings.ToUpper(postalCode)
	postalCode = strings.ReplaceAll(postalCode, " ", "")
	fmt.Println(postalCode)
	fmt.Println(DutchPostalCodeRegex.FindString(postalCode))
	fmt.Println(BelgianPostalCodeRegex.FindString(postalCode))
	if DutchPostalCodeRegex.FindString(postalCode) != "" && BelgianPostalCodeRegex.FindString(postalCode) != "" {
		return errors.New("postcode is onjuist. Geldige postcode voor Nederland is 1234AB, voor België 1234")
	}

	return nil
}

func validateDate(dateString string) error {
	_, err := time.Parse("2006-01-02", dateString)
	if err != nil {
		log.Println("Error parsing date:", err)
		return fmt.Errorf("de datum %s is niet correct", dateString)
	}
	return nil
}

// TODO make this function support non-dutch phone numbers too
func validatePhoneNumber(numberString string, label string) error {
	// assume the number is dutch in the first place
	_, err := phonenumbers.Parse(numberString, "NL")
	if err != nil {
		return fmt.Errorf("%s is niet correct", label)
	}

	return nil
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

	if iban == "" {
		return errors.New("IBAN-nummer is niet ingevuld")
	}

	resp, err := http.Get("https://openiban.com/validate/" + iban)
	if err != nil {
		return errors.New("kon IBAN niet valideren, probeer het later opnieuw")
	}
	var ibanval IBANValidationResponse

	err = json.NewDecoder(resp.Body).Decode(&ibanval)
	if err != nil {
		return errors.New("serverfout tijdens het valideren van de IBAN. Neem contact op met de vereniging")
	}

	if ibanval.Valid {
		return nil
	}

	return errors.New("IBAN is ongeldig: controleer of je alles goed hebt overgenomen")
}

func validateEmail(email string) error {
	_, err := mail.ParseAddress(email)
	if err != nil {
		return errors.New("email is ongeldig, probeer het zo: voorbeeld@svpromptusimperii.nl")
	}

	return nil
}
