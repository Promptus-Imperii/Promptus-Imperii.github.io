package main

import (
	"bytes"
	"encoding/base64"
	"io"
	"mime"
	"mime/multipart"
	"net/mail"
	"strings"
	"testing"
)

func testSignup(t *testing.T) Signup {
	t.Helper()
	s, errs := parseSignup(validForm(), testNow)
	if len(errs) != 0 {
		t.Fatal(errs)
	}
	return s
}

func TestSignupCSV(t *testing.T) {
	data, err := testSignup(t).csv()
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.HasPrefix(data, []byte("\xEF\xBB\xBF")) {
		t.Error("missing UTF-8 BOM")
	}
	want := "\xEF\xBB\xBF" +
		"Voornamen,Roepnaam,Tussenvoegsel,Achternaam,Geboortedatum,Adres,Postcode,Woonplaats,Land,E-mail,Opleiding,Type,Telefoonnummer,Cohortjaar,Noodnummer (Naam),Noodnummer (Telefoonnummer),IBAN,Naam rekeninghouder\n" +
		"Zoë Johanna,Zoë,van,Müller,2000-10-12,Lovensdijkstraat 61,4818 AJ,Breda,Nederland,zoe@example.org,,Lid,+31612345678,2025/2026,Greetje de Vries,+31687654321,NL91ABNA0417164300,Z. J. van Müller\n"
	if string(data) != want {
		t.Errorf("csv =\n%q\nwant\n%q", data, want)
	}
}

func TestFullName(t *testing.T) {
	s := testSignup(t)
	if s.fullName() != "Zoë van Müller" {
		t.Errorf("fullName = %q", s.fullName())
	}
	s.Nickname, s.Infix = "", ""
	if s.fullName() != "Zoë Johanna Müller" {
		t.Errorf("fullName = %q", s.fullName())
	}
}

func TestBuildMessage(t *testing.T) {
	csv := []byte("\xEF\xBB\xBFnaam\nZoë\n")
	gm, err := buildMessage("server@example.org", Mail{
		To:         "info@example.org",
		Subject:    "Nieuwe aanmelding lid: Zoë Müller",
		Body:       "Nieuw lid aangemeld, zie bijlage. Groetjes, Zoë",
		Attachment: &Attachment{Name: "nieuw_lid.csv", ContentType: "text/csv; charset=utf-8", Data: csv},
	})
	if err != nil {
		t.Fatal(err)
	}
	var raw bytes.Buffer
	if _, err := gm.WriteTo(&raw); err != nil {
		t.Fatal(err)
	}

	msg, err := mail.ReadMessage(&raw)
	if err != nil {
		t.Fatal(err)
	}
	subject, _ := new(mime.WordDecoder).DecodeHeader(msg.Header.Get("Subject"))
	if subject != "Nieuwe aanmelding lid: Zoë Müller" {
		t.Errorf("subject = %q", subject)
	}
	if !strings.Contains(msg.Header.Get("From"), "server@example.org") || !strings.Contains(msg.Header.Get("To"), "info@example.org") {
		t.Errorf("headers = %v", msg.Header)
	}

	_, params, _ := mime.ParseMediaType(msg.Header.Get("Content-Type"))
	parts := multipart.NewReader(msg.Body, params["boundary"])

	body, _ := parts.NextPart()
	text, _ := io.ReadAll(body) // multipart decodes quoted-printable
	if string(text) != "Nieuw lid aangemeld, zie bijlage. Groetjes, Zoë" {
		t.Errorf("body = %q", text)
	}

	att, _ := parts.NextPart()
	if att.FileName() != "nieuw_lid.csv" || !strings.HasPrefix(att.Header.Get("Content-Type"), "text/csv") {
		t.Errorf("attachment headers = %v", att.Header)
	}
	encoded, _ := io.ReadAll(att)
	decoded, err := base64.StdEncoding.DecodeString(strings.ReplaceAll(string(encoded), "\r\n", ""))
	if err != nil || !bytes.Equal(decoded, csv) {
		t.Errorf("attachment = %q (%v)", decoded, err)
	}
}
