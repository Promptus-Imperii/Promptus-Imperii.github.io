package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDevHandler(t *testing.T) {
	hugo := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.WriteString(w, "hugo "+r.URL.Path)
	}))
	defer hugo.Close()
	api := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.WriteString(w, "api "+r.URL.Path)
	})

	h, err := devHandler(api, hugo.URL)
	if err != nil {
		t.Fatal(err)
	}
	for path, want := range map[string]string{
		"/api/signup": "api /api/signup",
		"/contact/":   "hugo /contact/",
	} {
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, httptest.NewRequestWithContext(t.Context(), "GET", path, nil))
		if rec.Body.String() != want {
			t.Errorf("%s: %q", path, rec.Body.String())
		}
	}
}

func TestServerOutlastsTwoMails(t *testing.T) {
	srv := newServer(http.NotFoundHandler())
	if worst := 2 * smtpTimeout; srv.WriteTimeout <= worst {
		t.Errorf("WriteTimeout %v must exceed worst-case mail time %v", srv.WriteTimeout, worst)
	}
}
