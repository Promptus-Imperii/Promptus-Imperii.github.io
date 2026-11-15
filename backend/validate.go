package main

import (
	"fmt"
	"net/mail"
	"net/url"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"
)

const (
	maxFieldLen  = 200
	formulaStart = "=+-@"
)

var (
	nlPostalCode = regexp.MustCompile(`^[1-9][0-9]{3}[A-Z]{2}$`)
	bePostalCode = regexp.MustCompile(`^[1-9][0-9]{3}$`)
	// Any country: + and 8 to 15 digits, the international maximum.
	internationalPhone = regexp.MustCompile(`^\+[1-9][0-9]{7,14}$`)
	// Dutch numbers are +31 and 9 digits. Belgian ones are +32 and 8 (landline) or 9 (mobile) digits.
	phoneFormats = map[string]*regexp.Regexp{
		"+31": regexp.MustCompile(`^\+31[1-9][0-9]{8}$`),
		"+32": regexp.MustCompile(`^\+32[1-9][0-9]{7,8}$`),
	}

	cohortYear = regexp.MustCompile(`^([0-9]{4})/([0-9]{4})$`)
	ibanFormat = regexp.MustCompile(`^[A-Z]{2}[0-9]{2}[A-Z0-9]{11,30}$`)

	// "(0)" comes first: in "+31 (0)6 …" the 0 must go too.
	phoneSeparators = strings.NewReplacer("(0)", "", " ", "", "-", "", ".", "", "(", "", ")", "")
)

// rule describes one form field. The validate loop handles trimming, required, length,
// control characters and Excel formulas; check only covers what is specific to the field.
type rule struct {
	name, label string
	optional    bool
	// skipFormulaCheck is for fields whose check already allows only a strict format,
	// such as phone numbers, which can start with +.
	skipFormulaCheck bool
	normalize        func(string) string
	check            checkFunc
}

// checkFunc gets the normalized value and the values of the fields validated before it.
// It returns the value to keep, or an error message for the user.
type checkFunc func(value string, earlier map[string]string) (kept, msg string)

// validate runs the rules in order. It returns the cleaned values by name and Dutch error messages.
func validate(form url.Values, rules []rule) (map[string]string, []string) {
	values := map[string]string{}
	var errs []string
	for _, field := range rules {
		value := strings.TrimSpace(form.Get(field.name))
		if field.normalize != nil {
			value = field.normalize(value)
		}
		var msg string
		switch {
		case value == "" && field.optional:
		case value == "":
			msg = field.label + " is verplicht"
		case utf8.RuneCountInString(value) > maxFieldLen:
			msg = fmt.Sprintf("%s is te lang (maximaal %d tekens)", field.label, maxFieldLen)
		case strings.ContainsFunc(value, unicode.IsControl):
			msg = field.label + " bevat ongeldige tekens"
		case !field.skipFormulaCheck && strings.ContainsRune(formulaStart, rune(value[0])):
			// Excel would run the value as a formula when the secretary opens the CSV.
			msg = field.label + " mag niet beginnen met =, +, - of @"
		case field.check != nil:
			value, msg = field.check(value, values)
		}
		if msg != "" {
			errs = append(errs, msg)
			value = ""
		}
		values[field.name] = value
	}
	return values, errs
}

func upperNoSpaces(s string) string {
	return strings.ToUpper(strings.ReplaceAll(s, " ", ""))
}

func stripPhoneSeparators(s string) string {
	return phoneSeparators.Replace(s)
}

func checkOneOf(label string, options ...string) checkFunc {
	return func(value string, _ map[string]string) (kept, msg string) {
		if !slices.Contains(options, value) {
			return "", label + " is ongeldig"
		}
		return value, ""
	}
}

// checkPhone only accepts international numbers, such as +31612345678.
func checkPhone(label string) checkFunc {
	return func(value string, _ map[string]string) (kept, msg string) {
		format := internationalPhone
		if f, ok := phoneFormats[value[:min(3, len(value))]]; ok {
			format = f
		}
		if !format.MatchString(value) {
			return "", label + " is ongeldig. Begin met de landcode, zoals +31612345678 of +32467300512"
		}
		return value, ""
	}
}

func checkBirthDate(now time.Time) checkFunc {
	return func(value string, _ map[string]string) (kept, msg string) {
		date, err := time.Parse(time.DateOnly, value)
		if err != nil || !date.Before(now) || date.Year() < 1900 {
			return "", "Geboortedatum is ongeldig"
		}
		return value, ""
	}
}

// checkPostalCode depends on the country, so the country rule must come before it.
func checkPostalCode(value string, earlier map[string]string) (kept, msg string) {
	switch earlier["country"] {
	case "NL":
		if !nlPostalCode.MatchString(value) {
			return "", "Postcode is ongeldig. Een Nederlandse postcode ziet er zo uit: 1234 AB"
		}
		return value[:4] + " " + value[4:], ""
	case "BE":
		if !bePostalCode.MatchString(value) {
			return "", "Postcode is ongeldig. Een Belgische postcode ziet er zo uit: 1234"
		}
	}
	return value, ""
}

func checkEmail(value string, _ map[string]string) (kept, msg string) {
	if address, err := mail.ParseAddress(value); err != nil || address.Address != value {
		return "", "E-mailadres is ongeldig, probeer het zo: voorbeeld@svpromptusimperii.nl"
	}
	return value, ""
}

func checkCohortYear(value string, _ map[string]string) (kept, msg string) {
	const invalid = "Cohortjaar is ongeldig, gebruik het formaat 2025/2026"
	match := cohortYear.FindStringSubmatch(value)
	if match == nil {
		return "", invalid
	}
	first, _ := strconv.Atoi(match[1])
	second, _ := strconv.Atoi(match[2])
	if second != first+1 {
		return "", invalid
	}
	return value, ""
}

func checkIBAN(value string, _ map[string]string) (kept, msg string) {
	if !validIBAN(value) {
		return "", "IBAN is ongeldig: controleer of je alles goed hebt overgenomen"
	}
	return value, ""
}

// validIBAN checks the format and the mod-97 checksum.
func validIBAN(s string) bool {
	if !ibanFormat.MatchString(s) {
		return false
	}
	remainder := 0
	for _, c := range s[4:] + s[:4] {
		if c >= 'A' {
			remainder = (remainder*100 + int(c-'A'+10)) % 97
		} else {
			remainder = (remainder*10 + int(c-'0')) % 97
		}
	}
	return remainder == 1
}
