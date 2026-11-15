package main

import "net/http"

// handleEmailReveal shows the contact address once the captcha is solved, so scrapers can't read it.
func (a *app) handleEmailReveal(w http.ResponseWriter, r *http.Request) {
	if _, ok := verifyAltcha(a.cfg.AltchaKey, r.PostForm.Get("altcha")); !ok {
		writeErrors(w, http.StatusUnprocessableEntity, []string{msgCaptcha})
		return
	}
	writeHTML(w, http.StatusOK, render("email.html", a.cfg.ContactEmail))
}
