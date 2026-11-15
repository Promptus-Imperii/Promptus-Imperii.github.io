package main

import (
	"encoding/base64"
	"encoding/json"
	"os"
	"testing"
	"time"

	altcha "github.com/altcha-org/altcha-lib-go/v2"
)

var (
	testNow = time.Date(2026, 10, 5, 12, 0, 0, 0, time.UTC)
	testKey = []byte("test-key")
)

// TestMain makes challenges cheap to solve, so the tests stay fast.
func TestMain(m *testing.M) {
	altchaCost = 1
	randomAltchaCounter = func() int { return 3 }
	os.Exit(m.Run())
}

// solve answers a challenge like the Altcha widget does and returns the base64 payload.
func solve(t testing.TB, c altcha.Challenge) string {
	t.Helper()
	s, err := altcha.SolveChallenge(altcha.SolveChallengeOptions{Challenge: c, DeriveKey: altcha.DeriveKeyPBKDF2()})
	if err != nil || s == nil {
		t.Fatalf("no solution: %v", err)
	}
	b, _ := json.Marshal(altcha.Payload{Challenge: c, Solution: *s})
	return base64.StdEncoding.EncodeToString(b)
}

// challenge creates a challenge that is valid now (the library checks expiry against the real clock).
func challenge(t testing.TB) altcha.Challenge {
	t.Helper()
	c, err := newChallenge(testKey, time.Now())
	if err != nil {
		t.Fatal(err)
	}
	return c
}

func TestAltchaRoundTrip(t *testing.T) {
	c := challenge(t)
	if c.Parameters.Algorithm != "PBKDF2/SHA-256" || c.Parameters.ExpiresAt == 0 || c.Signature == "" {
		t.Fatalf("bad challenge %+v", c)
	}
	if !valid(testKey, solve(t, c)) {
		t.Error("valid payload rejected")
	}
}

func TestAltchaRejects(t *testing.T) {
	payload := solve(t, challenge(t))
	if valid([]byte("other-key"), payload) {
		t.Error("accepted with wrong key")
	}
	if valid(testKey, "not base64!") || valid(testKey, "") || valid(testKey, "e30=") {
		t.Error("accepted garbage")
	}

	expired, err := newChallenge(testKey, time.Now().Add(-altchaTTL-time.Minute))
	if err != nil {
		t.Fatal(err)
	}
	if valid(testKey, solve(t, expired)) {
		t.Error("accepted expired payload")
	}

	forged := challenge(t)
	forged.Signature = "00"
	if valid(testKey, solve(t, forged)) {
		t.Error("accepted forged signature")
	}

	var p altcha.Payload
	raw, _ := base64.StdEncoding.DecodeString(solve(t, challenge(t)))
	if err := json.Unmarshal(raw, &p); err != nil {
		t.Fatal(err)
	}
	p.Solution.Counter++
	b, _ := json.Marshal(p)
	if valid(testKey, base64.StdEncoding.EncodeToString(b)) {
		t.Error("accepted wrong solution")
	}
}

func valid(key []byte, payload string) bool {
	_, ok := verifyAltcha(key, payload)
	return ok
}

func TestUsedCaptchas(t *testing.T) {
	u := newUsedCaptchas()
	c := challenge(t)
	now := time.Now()
	if !u.claim(c, now) || u.claim(c, now) {
		t.Fatal("second claim must fail")
	}
	u.release(c)
	if !u.claim(c, now) {
		t.Error("released captcha must be claimable")
	}
	if !u.claim(c, now.Add(altchaTTL+time.Minute)) {
		t.Error("expired entry must be dropped")
	}
}
