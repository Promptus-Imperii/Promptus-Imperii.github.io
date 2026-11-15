// Command backend serves the website's signup form, contact address reveal and Altcha captcha.
package main

import (
	"errors"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"time"
)

type app struct {
	cfg      Config
	mailer   Mailer
	captchas *usedCaptchas
	now      func() time.Time
}

func newApp(cfg Config) *app {
	var m Mailer = smtpMailer{host: cfg.SMTPHost, port: cfg.SMTPPort, from: cfg.SMTPFrom}
	if cfg.Dev {
		m = logMailer{}
	}
	return &app{cfg: cfg, mailer: m, captchas: newUsedCaptchas(), now: time.Now}
}

// routes accepts paths with or without /api, in case the proxy strips it.
func (a *app) routes() http.Handler {
	mux := http.NewServeMux()
	for _, prefix := range []string{"", "/api"} {
		mux.HandleFunc("GET "+prefix+"/captcha-challenge", a.handleChallenge)
		mux.HandleFunc("POST "+prefix+"/signup", parseForm(a.handleSignup))
		mux.HandleFunc("POST "+prefix+"/email", parseForm(a.handleEmailReveal))
	}
	return mux
}

const (
	msgCaptcha = "Een geldige captcha is vereist. Probeer de pagina te herladen (je formuliervelden blijven bestaan)"
	msgBadForm = "Het formulier kon niet worden verwerkt. Probeer het opnieuw."
)

// parseForm caps the body at 64 KB and parses it, so handlers can read r.PostForm.
func parseForm(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, 64<<10)
		if err := r.ParseForm(); err != nil {
			status := http.StatusBadRequest
			var tooLarge *http.MaxBytesError
			if errors.As(err, &tooLarge) {
				status = http.StatusRequestEntityTooLarge
			}
			writeErrors(w, status, []string{msgBadForm})
			return
		}
		next(w, r)
	}
}

func main() {
	cfg, err := loadConfig(os.Getenv)
	if err != nil {
		log.Fatal(err)
	}
	handler := newApp(cfg).routes()
	if cfg.Dev {
		if handler, err = devHandler(handler, cfg.DevProxy); err != nil {
			log.Fatal(err)
		}
	}
	log.Printf("listening on :3000 (dev: %v)", cfg.Dev)
	log.Fatal(newServer(handler).ListenAndServe())
}

// newServer's write timeout fits two slow mail sends.
func newServer(handler http.Handler) *http.Server {
	return &http.Server{
		Addr:              ":3000",
		Handler:           handler,
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      2*smtpTimeout + 20*time.Second,
	}
}

// devHandler proxies everything except /api/ to `hugo server`, so dev is same-origin.
func devHandler(api http.Handler, target string) (http.Handler, error) {
	u, err := url.Parse(target)
	if err != nil {
		return nil, err
	}
	proxy := httputil.NewSingleHostReverseProxy(u)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/api/") {
			api.ServeHTTP(w, r)
			return
		}
		proxy.ServeHTTP(w, r)
	}), nil
}
