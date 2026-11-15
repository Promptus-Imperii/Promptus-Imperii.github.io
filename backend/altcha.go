package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"log"
	"math/big"
	"net/http"
	"sync"
	"time"

	altcha "github.com/altcha-org/altcha-lib-go/v2"
)

// A solved challenge stays valid until it expires, so users can fix form errors without solving it again.
const altchaTTL = 20 * time.Minute

// About 2M PBKDF2 iterations on average, roughly a second on a phone. Tests lower these.
var (
	altchaCost = 1000
	// The counter is what the browser has to find, so it must not be predictable.
	randomAltchaCounter = func() int {
		n, _ := rand.Int(rand.Reader, big.NewInt(2000)) // never fails: crypto/rand panics instead
		return 1000 + int(n.Int64())
	}
)

func newChallenge(key []byte, now time.Time) (altcha.Challenge, error) {
	expires := now.Add(altchaTTL)
	counter := randomAltchaCounter()
	return altcha.CreateChallenge(altcha.CreateChallengeOptions{
		Algorithm:           "PBKDF2/SHA-256", // required, the library has no default
		DeriveKey:           altcha.DeriveKeyPBKDF2(),
		HMACSignatureSecret: string(key),
		Cost:                altchaCost,
		Counter:             &counter,
		ExpiresAt:           &expires,
	})
}

// verifyAltcha checks the widget payload and returns its challenge. Expiry uses the real clock.
func verifyAltcha(key []byte, payload string) (altcha.Challenge, bool) {
	raw, err := base64.StdEncoding.DecodeString(payload)
	if err != nil {
		return altcha.Challenge{}, false
	}
	var p altcha.Payload
	if json.Unmarshal(raw, &p) != nil {
		return altcha.Challenge{}, false
	}
	result, err := altcha.VerifySolution(altcha.VerifySolutionOptions{
		Challenge:           p.Challenge,
		Solution:            p.Solution,
		DeriveKey:           altcha.DeriveKeyPBKDF2(),
		HMACSignatureSecret: string(key),
	})
	return p.Challenge, err == nil && result.Verified
}

// usedCaptchas lets a solved captcha make only one signup. Users can still retry after
// validation errors, because a captcha is only claimed once the form is valid.
// Entries are dropped after they expire; a restart forgets them, which is fine for 20 minutes.
type usedCaptchas struct {
	mu      sync.Mutex
	expires map[string]time.Time // challenge signature -> expiry
}

func newUsedCaptchas() *usedCaptchas {
	return &usedCaptchas{expires: map[string]time.Time{}}
}

// claim marks the challenge as used. It returns false if it was used before.
func (u *usedCaptchas) claim(c altcha.Challenge, now time.Time) bool {
	u.mu.Lock()
	defer u.mu.Unlock()
	for sig, exp := range u.expires {
		if now.After(exp) {
			delete(u.expires, sig)
		}
	}
	if _, used := u.expires[c.Signature]; used {
		return false
	}
	u.expires[c.Signature] = time.Unix(c.Parameters.ExpiresAt, 0)
	return true
}

// release makes the challenge usable again, for when the signup failed on our side.
func (u *usedCaptchas) release(c altcha.Challenge) {
	u.mu.Lock()
	defer u.mu.Unlock()
	delete(u.expires, c.Signature)
}

func (a *app) handleChallenge(w http.ResponseWriter, _ *http.Request) {
	c, err := newChallenge(a.cfg.AltchaKey, a.now())
	if err != nil {
		log.Printf("altcha challenge: %v", err)
		http.Error(w, "captcha unavailable", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	_ = json.NewEncoder(w).Encode(c)
}
