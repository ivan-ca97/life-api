package service

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/ivan-ca97/life/internal/features/auth/domain"
	"github.com/ivan-ca97/life/internal/features/auth/ports"
)

type googleTokenVerifier struct {
	httpClient *http.Client
}

var _ ports.GoogleTokenVerifier = (*googleTokenVerifier)(nil)

func NewGoogleTokenVerifier() *googleTokenVerifier {
	return &googleTokenVerifier{
		httpClient: &http.Client{},
	}
}

type tokenInfoResponse struct {
	Subject       string `json:"sub"`
	Email         string `json:"email"`
	EmailVerified string `json:"email_verified"`
	Audience      string `json:"aud"`
}

func (v *googleTokenVerifier) Verify(idToken, clientId string) (*ports.GoogleClaims, error) {
	tokenInfoURL := "https://oauth2.googleapis.com/tokeninfo?id_token=" + url.QueryEscape(idToken)
	response, err := v.httpClient.Get(tokenInfoURL)
	if err != nil {
		return nil, domain.ErrInvalidGoogleToken
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, domain.ErrInvalidGoogleToken
	}

	var tokenInfo tokenInfoResponse
	err = json.NewDecoder(response.Body).Decode(&tokenInfo)
	if err != nil {
		return nil, domain.ErrInvalidGoogleToken
	}

	if tokenInfo.Audience != clientId {
		return nil, domain.ErrInvalidGoogleToken
	}

	if tokenInfo.EmailVerified != "true" {
		return nil, domain.ErrInvalidGoogleToken
	}

	claims := &ports.GoogleClaims{
		Subject: tokenInfo.Subject,
		Email:   tokenInfo.Email,
	}
	return claims, nil
}
