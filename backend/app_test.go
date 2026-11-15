package main

import (
	"context"
	"encoding/json"
	"errors"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	altcha "github.com/altcha-org/altcha-lib-go/v2"
)

type fakeMailer struct {
	sent []Mail
	fail map[string]bool // recipient -> fail
}

func (f *fakeMailer) Send(_ context.Context, m Mail) error {
	if f.fail[m.To] {
		return errors.New("smtp down")
	}
	f.sent = append(f.sent, m)
	return nil
}

func testApp() (*app, *fakeMailer) {
	m := &fakeMailer{fail: map[string]bool{}}
	return &app{
		cfg:      Config{ContactEmail: "info@example.org", AltchaKey: testKey},
		mailer:   m,
		captchas: newUsedCaptchas(),
		now:      func() time.Time { return testNow },
	}, m
}

// post sends a form and checks the status code. It returns the body and the HX-Redirect header.
func post(t *testing.T, a *app, path string, form url.Values, wantCode int) (body, redirect string) {
	t.Helper()
	req := httptest.NewRequestWithContext(t.Context(), "POST", path, strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	a.routes().ServeHTTP(rec, req)
	if rec.Code != wantCode {
		t.Fatalf("%s: code %d, want %d, body %s", path, rec.Code, wantCode, rec.Body)
	}
	if rec.Body.Len() > 0 && rec.Header().Get("Content-Type") != "text/html; charset=utf-8" {
		t.Errorf("%s: content-type %q", path, rec.Header().Get("Content-Type"))
	}
	return rec.Body.String(), rec.Header().Get("HX-Redirect")
}

func withCaptcha(t *testing.T, f url.Values) url.Values {
	f.Set("altcha", solve(t, challenge(t)))
	return f
}

func TestChallengeEndpoint(t *testing.T) {
	a, _ := testApp()
	rec := httptest.NewRecorder()
	a.routes().ServeHTTP(rec, httptest.NewRequestWithContext(t.Context(), "GET", "/api/captcha-challenge", nil))
	var c altcha.Challenge
	if err := json.NewDecoder(rec.Body).Decode(&c); err != nil || c.Signature == "" {
		t.Fatalf("challenge %+v, %v", c, err)
	}
}

func TestSignupSuccess(t *testing.T) {
	a, m := testApp()
	_, redirect := post(t, a, "/api/signup", withCaptcha(t, validForm()), 204)
	if redirect != "/signup_success/" {
		t.Errorf("redirect = %q", redirect)
	}
	if len(m.sent) != 2 {
		t.Fatalf("sent %d mails", len(m.sent))
	}
	secretary, member := m.sent[0], m.sent[1]
	if secretary.To != "info@example.org" || secretary.Subject != "[Server] Nieuwe aanmelding lid: Zoë van Müller" || secretary.Attachment == nil {
		t.Errorf("secretary mail = %+v", secretary)
	}
	if member.To != "zoe@example.org" || !strings.Contains(member.Body, "info@example.org") {
		t.Errorf("member mail = %+v", member)
	}
}

func TestSignupValidationErrors(t *testing.T) {
	a, m := testApp()
	f := validForm()
	f.Set("iban", "NL00")
	body, _ := post(t, a, "/api/signup", withCaptcha(t, f), 422)
	if !strings.Contains(body, `<div id="errors"`) || !strings.Contains(body, "IBAN is ongeldig") {
		t.Errorf("body = %s", body)
	}
	if len(m.sent) != 0 {
		t.Error("mail sent despite errors")
	}
}

func TestSignupRetryWithSameCaptcha(t *testing.T) {
	a, m := testApp()
	f := withCaptcha(t, validForm())
	f.Set("iban", "NL00")
	post(t, a, "/api/signup", f, 422)
	f.Set("iban", "NL91ABNA0417164300")
	_, redirect := post(t, a, "/signup", f, 204)
	if redirect == "" || len(m.sent) != 2 {
		t.Errorf("retry failed: redirect %q, sent %d", redirect, len(m.sent))
	}
}

func TestSignupInvalidCaptcha(t *testing.T) {
	a, m := testApp()
	f := validForm()
	f.Set("altcha", "bogus")
	body, _ := post(t, a, "/api/signup", f, 422)
	if !strings.Contains(body, "Een geldige captcha is vereist") || len(m.sent) != 0 {
		t.Errorf("body = %s", body)
	}
}

func TestSignupSecretaryMailFailure(t *testing.T) {
	a, m := testApp()
	m.fail["info@example.org"] = true
	body, redirect := post(t, a, "/api/signup", withCaptcha(t, validForm()), 500)
	if !strings.Contains(body, "Er is iets fout gegaan") || !strings.Contains(body, "info@example.org") {
		t.Errorf("body = %s", body)
	}
	if redirect != "" || len(m.sent) != 0 {
		t.Error("must not succeed or send the confirmation")
	}
}

func TestSignupConfirmationFailureStillSucceeds(t *testing.T) {
	a, m := testApp()
	m.fail["zoe@example.org"] = true
	_, redirect := post(t, a, "/api/signup", withCaptcha(t, validForm()), 204)
	if redirect == "" || len(m.sent) != 1 {
		t.Errorf("redirect %q, sent %d", redirect, len(m.sent))
	}
}

func TestEmailReveal(t *testing.T) {
	a, _ := testApp()
	body, _ := post(t, a, "/api/email", withCaptcha(t, url.Values{}), 200)
	if !strings.Contains(body, `<div id="email-block"`) || !strings.Contains(body, `href="mailto:info@example.org"`) {
		t.Errorf("body = %s", body)
	}

	bad, _ := post(t, a, "/api/email", url.Values{"altcha": {"x"}}, 422)
	if strings.Contains(bad, "mailto:") || !strings.Contains(bad, "captcha") {
		t.Errorf("bad captcha body = %s", bad)
	}
}

func TestRoutesWithAndWithoutPrefix(t *testing.T) {
	a, _ := testApp()
	for _, p := range []string{"/api/email", "/email"} {
		if body, _ := post(t, a, p, withCaptcha(t, url.Values{}), 200); !strings.Contains(body, "mailto:") {
			t.Errorf("%s: %s", p, body)
		}
	}
	rec := httptest.NewRecorder()
	a.routes().ServeHTTP(rec, httptest.NewRequestWithContext(t.Context(), "GET", "/api/signup", nil))
	if rec.Code != 405 {
		t.Errorf("GET signup = %d", rec.Code)
	}
}

func TestSignupSubjectCutsOffLongNames(t *testing.T) {
	a, m := testApp()
	f := validForm()
	long := strings.Repeat("ë", 200)
	f.Set("nickname", long)
	f.Set("infix", long)
	f.Set("surname", long)
	post(t, a, "/api/signup", withCaptcha(t, f), 204)
	if len(m.sent) != 2 {
		t.Fatalf("sent %d mails", len(m.sent))
	}
	want := "[Server] Nieuwe aanmelding lid: " + strings.Repeat("ë", 60) + "…"
	if m.sent[0].Subject != want {
		t.Errorf("subject = %q", m.sent[0].Subject)
	}
	csv := string(m.sent[0].Attachment.Data)
	if !strings.Contains(csv, long+","+long+",") {
		t.Error("CSV must keep the full name")
	}
}

func TestSignupCaptchaOnlyOnce(t *testing.T) {
	a, m := testApp()
	f := withCaptcha(t, validForm())
	post(t, a, "/api/signup", f, 204)
	f.Set("email", "someone-else@example.org")
	if body, _ := post(t, a, "/api/signup", f, 422); !strings.Contains(body, "captcha") {
		t.Errorf("body = %s", body)
	}
	if len(m.sent) != 2 {
		t.Errorf("sent %d mails, want 2", len(m.sent))
	}
}

func TestSignupFailedMailReleasesCaptcha(t *testing.T) {
	a, m := testApp()
	m.fail["info@example.org"] = true
	f := withCaptcha(t, validForm())
	post(t, a, "/api/signup", f, 500)
	delete(m.fail, "info@example.org")
	post(t, a, "/api/signup", f, 204)
}

func TestSignupIgnoresQueryParameters(t *testing.T) {
	a, _ := testApp()
	f := withCaptcha(t, validForm())
	iban := f.Get("iban")
	f.Del("iban")
	if body, _ := post(t, a, "/api/signup?iban="+url.QueryEscape(iban), f, 422); !strings.Contains(body, "IBAN is verplicht") {
		t.Errorf("body = %s", body)
	}
}

func TestOversizedFormIsNotACaptchaError(t *testing.T) {
	a, _ := testApp()
	body, _ := post(t, a, "/api/signup", url.Values{"x": {strings.Repeat("a", 70<<10)}}, 413)
	if strings.Contains(body, "captcha") || !strings.Contains(body, "kon niet worden verwerkt") {
		t.Errorf("body = %s", body)
	}
}
