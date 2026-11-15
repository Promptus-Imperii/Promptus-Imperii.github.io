package main

import "net/http"

// Responses are HTML fragments for htmx. Each fragment is wrapped in <hx-partial>, which tells
// htmx where it goes, so the forms only need hx-post. The site tells htmx not to swap 4xx and 5xx
// responses (a proxy error page must not end up in the page), but partials are still placed,
// so our own errors show up with a 422 or 500 status.

func errorsHTML(msgs []string) string {
	return render("errors.html", msgs)
}

func writeHTML(w http.ResponseWriter, status int, html string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(status)
	_, _ = w.Write([]byte(html)) //nolint:gosec // html comes from html/template, which escapes user input
}

func writeErrors(w http.ResponseWriter, status int, msgs []string) {
	writeHTML(w, status, errorsHTML(msgs))
}

// redirect makes htmx do a full page navigation.
func redirect(w http.ResponseWriter, url string) {
	w.Header().Set("HX-Redirect", url)
	w.WriteHeader(http.StatusNoContent)
}
