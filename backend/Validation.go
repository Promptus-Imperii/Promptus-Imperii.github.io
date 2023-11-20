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
	// implemented to https://nl.wikipedia.org/wiki/Postcode#Postcodes_in_Nederland .
	// lookahead is not currently supported in standard library regex.
	PostalCodeRegex *regexp.Regexp = regexp.MustCompile(`^[1-9]\d{3}\w{2}$`)
	// this does not need to be fancy, just needs to check if it is somewhat valid.
	//
	// checks for `*@*.*`
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
		return errors.New("postcode is onjuist, probeer het zo: 4818 AJ")
	}

	if code[len(code)-2] == 'S' && strings.ContainsAny(string(code[len(code)-1]), "ADS") {
		return errors.New("onjuiste postcode")
	}

	return nil
}

// TODO make this function support non-dutch phone numbers too
func validatePhoneNumber(numberString string) error {
	// assume the number is dutch in the first place
	number, err := phonenumbers.Parse(numberString, "NL")
	if err != nil {
		return errors.New("dit is geen telefoonnummer")
	}

	if !phonenumbers.IsValidNumberForRegion(number, "NL") {
		return errors.New("geen geldig Nederlands nummer")
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
		return errors.New("geen IBAN gevonden")
	}

	resp, err := http.Get("https://openiban.com/validate/" + iban)
	if err != nil {
		return errors.New("kon IBAN niet valideren, probeer het later opnieuw")
	}
	var ibanval IBANValidationResponse

	err = json.NewDecoder(resp.Body).Decode(&ibanval)
	if err != nil {
		return errors.New("kon IBAN niet valideren (fout bij externe service)")
	}

	if ibanval.Valid {
		return nil
	}

	return errors.New("IBAN is ongeldig: controleer of je alles goed hebt overgenomen")
}

// DISCUSS: should this send the email already?
func validateEmail(email string) error {
	if EmailRegex.FindString(email) == "" {
		return errors.New("email is ongeldig")
	}

	// idk send the verification email?
	// definitely use a goroutine for that
	//
	// go sendVerificationEmail(email)

	return nil
}
