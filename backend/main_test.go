package main

import (
	"net/http"
	"strings"
	"testing"

	"github.com/gavv/httpexpect/v2"
	"github.com/gin-gonic/gin"
)

var correctUser = map[string]interface{}{
	"legal_first_names":              "boben b",
	"nickname":                       "bob",
	"infix":                          "de",
	"surname":                        "tak",
	"date_of_birth":                  "2024-03-23",
	"address":                        "Lovensdijkstaat 16",
	"postal_code":                    "4793AB",
	"city":                           "Breda",
	"country":                        "NL",
	"email":                          "jandevries@example.org",
	"education":                      "TI",
	"phone":                          "+31612345678",
	"cohort_year":                    "2022/2023",
	"emergency_contact_first_name":   "greetje",
	"emergency_contact_infix":        "de",
	"emergency_contact_surname":      "vries",
	"emergency_contact_phone_number": "+31687654321",
	"iban":                           "NL18RABO0123459876",
	"account_holder":                 "B. B. de Tak",
	"accept_contribution":            "on",
	"accept_terms_and_conditions":    "on",
}

func getGinHandler(t *testing.T) *httpexpect.Expect {
	// Create new gin instance
	handler := initRouter()
	// Create httpexpect instance
	gin.SetMode(gin.TestMode)
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
	e.POST("/api/signup").
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
	e.POST("/api/signup").
		WithJSON(userWithIncorrectPostalcodeNumbers).
		Expect().
		Status(http.StatusBadRequest).JSON().Object().
		HasValue("Errors", []string{"postcode is onjuist. Geldige postcode voor Nederland is 1234AB, voor België 1234"})

	e.POST("/api/signup").
		WithJSON(userWithIncorrectPostalcodeLetters).
		Expect().
		Status(http.StatusBadRequest).JSON().Object().
		HasValue("Errors", []string{"postcode is onjuist. Geldige postcode voor Nederland is 1234AB, voor België 1234"})
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

func TestEmailValidationRejectsInvalidEmail(t *testing.T) {
	err := validateEmail("@svpromptusimperii.nl")
	if err == nil {
		t.FailNow()
	}
}

func TestCohortYearValidationAcceptsValidCohortYear(t *testing.T) {
	err := validateCohortYear("2023/2024")
	if err != nil {
		t.FailNow()
	}
}

func TestCohortYearValidationRejectsInalidCohortYear(t *testing.T) {
	err := validateCohortYear("23/24")
	if err == nil {
		t.FailNow()
	}
}

func TestValidatePostalCodeWithValidDutchPostalCode1(t *testing.T) {
	_, err := validatePostalCode("1234AB")
	if err != nil {
		t.FailNow()
	}
}

func TestValidatePostalCodeWithValidDutchPostalCode2(t *testing.T) {
	_, err := validatePostalCode("1234 AB")
	if err != nil {
		t.FailNow()
	}
}

func TestValidatePostalCodeWithCorrectBelgianPostalCode(t *testing.T) {
	_, err := validatePostalCode("1234")
	if err != nil {
		t.FailNow()
	}
}

func TestValidateDateWithValidDateFormat(t *testing.T) {
	err := validateDate("2024-03-23")
	if err != nil {
		t.FailNow()
	}
}

func TestValidateDateWithInvalidDateFormat(t *testing.T) {
	err := validateDate("23/03/2024")
	if err == nil {
		t.FailNow()
	}
}

func TestValidatePhoneNumberWithValidDutchPhoneNumber(t *testing.T) {
	err := validatePhoneNumber("+31612345678", "Phone")
	if err != nil {
		t.FailNow()
	}
}

func TestValidatePhoneNumberWithInvalidDutchPhoneNumberFormat(t *testing.T) {
	err := validatePhoneNumber("12345678", "Phone")
	if err == nil {
		t.FailNow()
	}
}
func TestValidatePhoneNumberWithInvalidDutchPhoneNumberFormat2(t *testing.T) {
	err := validatePhoneNumber("0612345678", "Phone")
	if err == nil {
		t.FailNow()
	}
}
func TestValidatePhoneNumberWithValidBelgianPhoneNumber(t *testing.T) {
	err := validatePhoneNumber("+32466117160", "Phone")
	if err != nil {
		t.FailNow()
	}
}

func TestValidateCohortYearWithValidCohortYearFormat(t *testing.T) {
	err := validateCohortYear("2023/2024")
	if err != nil {
		t.FailNow()
	}
}

func TestValidateCohortYearWithInvalidCohortYearFormat(t *testing.T) {
	err := validateCohortYear("23/24")
	if err == nil {
		t.FailNow()
	}
}

func TestIbanValidationAcceptsIbanWithSpaces(t *testing.T) {
	err := validateIBAN("NL12 ABNA 8803 9263 72")
	if err != nil {
		t.FailNow()
	}
}

func TestIbanValidationRejectsIbanWithUrlCharacters(t *testing.T) {
	err := validateIBAN("NL12ABNA8803926372/../../?x=1")
	if err == nil {
		t.FailNow()
	}
}

func TestSignupShouldRejectOversizedBody(t *testing.T) {
	// Arrange
	e := getGinHandler(t)
	oversizedUser := map[string]interface{}{"legal_first_names": strings.Repeat("a", MAX_REQUEST_BODY_BYTES)}

	// Act & Assert
	e.POST("/api/signup").
		WithJSON(oversizedUser).
		Expect().
		Status(http.StatusBadRequest)
}

func TestEscapeCSVFormulaEscapesFormulas(t *testing.T) {
	for _, value := range []string{"=HYPERLINK(\"https://example.org\")", "+SUM(A1)", "-1+1", "@cmd", "\t=1"} {
		if escapeCSVFormula(value) != "'"+value {
			t.Fatalf("%q was not escaped", value)
		}
	}
}

func TestEscapeCSVFormulaKeepsNormalValues(t *testing.T) {
	for _, value := range []string{"", "Jan", "+31612345678", "+32 466 11 71 60", "Lovensdijkstraat 16"} {
		if escapeCSVFormula(value) != value {
			t.Fatalf("%q should not be escaped", value)
		}
	}
}

func TestExportConvertsBelgium(t *testing.T) {
	member := PISignUp{Country: "BE"}
	if member.ToPISignUpExport().Land != "België" {
		t.FailNow()
	}
}

func TestExportUsesEmergencyContactInfix(t *testing.T) {
	member := PISignUp{EmergencyContactFirstName: "Greetje", EmergencyContactInfix: "de", EmergencyContactSurname: "Vries"}
	if member.ToPISignUpExport().Noodnummer_Naam != "Greetje de Vries" {
		t.FailNow()
	}
}
