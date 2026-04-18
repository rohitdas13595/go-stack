package oauth

import (
	"fmt"
	"net/http"
)

// Profile is returned by providers after OAuth success.
type Profile struct {
	Email string
	Name  string
	ID    string
}

// Provider configures an OAuth2 provider (stub: wire golang.org/x/oauth2 in app).
type Provider struct {
	Name        string
	AuthURL     string
	TokenURL    string
	ClientID    string
	Secret      string
	RedirectURL string
}

// Google returns stub Google provider config.
func Google(clientID, secret string) Provider {
	return Provider{
		Name:        "google",
		AuthURL:     "https://accounts.google.com/o/oauth2/v2/auth",
		TokenURL:    "https://oauth2.googleapis.com/token",
		ClientID:    clientID,
		Secret:      secret,
		RedirectURL: "/auth/google/callback",
	}
}

// GitHub returns stub GitHub provider config.
func GitHub(clientID, secret string) Provider {
	return Provider{
		Name:        "github",
		AuthURL:     "https://github.com/login/oauth/authorize",
		TokenURL:    "https://github.com/login/oauth/access_token",
		ClientID:    clientID,
		Secret:      secret,
		RedirectURL: "/auth/github/callback",
	}
}

// StartRedirect redirects to provider authorize URL (caller must build query).
func StartRedirect(w http.ResponseWriter, r *http.Request, p Provider) error {
	_ = p
	return fmt.Errorf("oauth: implement authorize redirect with state + PKCE")
}
