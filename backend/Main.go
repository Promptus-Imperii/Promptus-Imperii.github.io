package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nyaruka/phonenumbers"
)

func main() {
	router := gin.Default()

	// register end point
	router.POST("/signup", handleSignUp)

	router.Run(":8080")
}

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

func handleSignUp(context *gin.Context) {
	var signup PISignUP

	err := json.NewDecoder(context.Request.Body).Decode(&signup)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO(jonas): make sure that if a validation fails,
	// a descriptive error is returned so a clear error message can be
	// displayed on the page.

	// normalize to save some time on regex :D
	signup.PostalCode = strings.ReplaceAll(signup.PostalCode, " ", "")

	// oh boy i love validating
	err = validatePostalCode(signup.PostalCode)
	if err != nil {
		returnErr(context, err)
		return
	}

	// validate own phone number
	err = validatePhoneNumber(signup.Member.PhoneNumber)
	if err != nil {
		returnErr(context, err)
		return
	}

	err = validateIBAN(signup.IBAN)
	if err != nil {
		returnErr(context, err)
		return
	}

	err = validatePhoneNumber(signup.EmergencyContact.PhoneNumber)
	if err != nil {
		returnErr(context, errors.New("noodcontact: "+err.Error()))
		return
	}

	err = validateEmail(signup.Email)
	if err != nil {
		returnErr(context, err)
		return
	}

	// at this point everything *should* be okay
	// sending the message already might be early

	context.JSON(http.StatusOK, gin.H{"success": "aangemeld! welkom bij S.V. Promptus Imperii!"})
}

func returnErr(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
}

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
		return errors.New("onjuiste postcode, probeer het zo: 4818 AJ")
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
		return errors.New("kon IBAN niet valideren")
	}
	var ibanval IBANValidationResponse

	err = json.NewDecoder(resp.Body).Decode(&ibanval)
	if err != nil {
		return errors.New("kon IBAN niet valideren (fout bij externe service)")
	}

	if ibanval.Valid {
		return nil
	}

	return errors.New("iban is ongeldig, controleer of je alles goed hebt overgenomen")
}

// DISCUSS: should this send the email already?
func validateEmail(email string) error {
	if EmailRegex.FindString(email) == "" {
		return errors.New("lijkt niet op een email")
	}

	// idk send the verification email?
	// definitely use a goroutine for that
	//
	// go sendVerificationEmail(email)

	return nil
}
