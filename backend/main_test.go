package main

import (
	"net/http"
	"testing"

	"github.com/gavv/httpexpect/v2"
	"github.com/gin-gonic/gin"
)

// Rid of debug output
func init() {
	gin.SetMode(gin.TestMode)
}

var correctUser = map[string]interface{}{
	"legalfirstnames": "boben b",
	"member": map[string]interface{}{
		"firstname": "bob",
		"infix":     "de",
		"lastname":  "tak",
		"phone":     "+31612345678",
	},
	"date_of_birth": "2000-10-12T00:00:00Z",
	"address":       "Lovensdijkstaat 16",
	"postal_code":   "4793RR",
	"city":          "Breda",
	"email":         "jandevries@example.org",
	"course":        "TI",
	"cohort":        "2022/2023",
	"emergency_contact": map[string]interface{}{
		"firstname": "greetje",
		"infix":     "de",
		"lastname":  "vries",
		"phone":     "+31687654321",
	},
	"iban":           "NL18RABO0123459876",
	"account_holder": "B. B. de Tak",
}

func getGinHandler(t *testing.T) *httpexpect.Expect {
	// Create new gin instance
	handler := initRouter()
	// Create httpexpect instance
	return httpexpect.WithConfig(httpexpect.Config{
		Client: &http.Client{
			Transport: httpexpect.NewBinder(handler),
			Jar:       httpexpect.NewCookieJar(),
		},
		Reporter: httpexpect.NewAssertReporter(t),
		Printers: []httpexpect.Printer{
			httpexpect.NewDebugPrinter(t, true),
		},
	})
}

func TestSignupShouldReturnSuccessWhenUserIsCorrect(t *testing.T) {
	// Arrange
	e := getGinHandler(t)

	// Act & Assert
	e.POST("/signup").
		WithJSON(correctUser).
		Expect().
		Status(http.StatusOK).JSON().
		Object().HasValue("Success", "Registration successful.")
}

func TestSignupShouldReturnErrorWhenPostalCodeIsInvalid(t *testing.T) {
	// Arrange
	e := getGinHandler(t)
	userWithIncorrectPostalcodeNumbers := correctUser
	userWithIncorrectPostalcodeLetters := correctUser

	userWithIncorrectPostalcodeNumbers["postal_code"] = "132NV"
	userWithIncorrectPostalcodeLetters["postal_code"] = "1323N"

	// Act & Assert
	e.POST("/signup").
		WithJSON(userWithIncorrectPostalcodeNumbers).
		Expect().
		Status(http.StatusBadRequest).JSON().
		Object().HasValue("Error", "postcode is onjuist, probeer het zo: 4818 AJ")

	e.POST("/signup").
		WithJSON(userWithIncorrectPostalcodeLetters).
		Expect().
		Status(http.StatusBadRequest).JSON().
		Object().HasValue("Error", "postcode is onjuist, probeer het zo: 4818 AJ")
}

func TestIbanValidationAcceptsValidIban(t *testing.T) {
	err := validateIBAN("NL12ABNA8803926372")
	if err != nil {
		t.FailNow()
	}
}

func TestIbanValidationRejectsEmptyIban(t *testing.T) {
	err := validateIBAN("")
	if err == nil {
		t.FailNow()
	}
}

func TestIbanValidationRejectsImproperIban(t *testing.T) {
	err := validateIBAN("NL12ABNA88039263")
	if err == nil {
		t.FailNow()
	}
}

func TestEmailValidationAcceptsValidEmail(t *testing.T) {
	err := validateEmail("hello@svpromptusimperii.nl")
	if err != nil {
		t.FailNow()
	}
}

func TestEmailValidationRejectsInalidEmail(t *testing.T) {
	err := validateEmail("@svpromptusimperii.nl")
	if err == nil {
		t.FailNow()
	}
}
