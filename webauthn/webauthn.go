package webauthn

import (
	"fmt"
	"net/http"
)

// Config holds RP settings for go-webauthn integration (stub registration).
type Config struct {
	RPDisplayName string
	RPID          string
	Origins       []string
}

// SessionData placeholder until full WebAuthn flow is wired.
type SessionData struct{}

// BeginRegistration returns a stub response (implement with github.com/go-webauthn/webauthn).
func BeginRegistration(w http.ResponseWriter, r *http.Request, cfg Config) error {
	_ = cfg
	return fmt.Errorf("webauthn: BeginRegistration not wired; add go-webauthn dependency in app")
}

// FinishRegistration completes registration (stub).
func FinishRegistration(w http.ResponseWriter, r *http.Request, cfg Config) error {
	_ = cfg
	return fmt.Errorf("webauthn: FinishRegistration not wired")
}
