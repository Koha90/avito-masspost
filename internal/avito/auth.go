package avito

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Auth provides access to Avito authentication methods.
type Auth struct {
	client       *Client
	clientID     string
	clientSecret string
}

// Token contains OAuth token data.
type Token struct {
	AccessToken string
	TokenType   string
	ExpiresAt   time.Time
}

type tokenDTO struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int64  `json:"expires_in"`
}

// NewAuth returns a new authentication client.
func NewAuth(client *Client, clientID, clientSecret string) *Auth {
	return &Auth{
		client:       client,
		clientID:     clientID,
		clientSecret: clientSecret,
	}
}

// Token requests a new access token.
func (a *Auth) Token(ctx context.Context) (Token, error) {
	form := url.Values{}
	form.Set("grant_type", "client_credentials")
	form.Set("client_id", a.clientID)
	form.Set("client_secret", a.clientSecret)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		strings.TrimRight(a.client.baseURL, "/")+"/token",
		strings.NewReader(form.Encode()),
	)
	if err != nil {
		return Token{}, fmt.Errorf("build token request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := a.client.httpClient.Do(req)
	if err != nil {
		return Token{}, fmt.Errorf("perform token request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusMultipleChoices {
		return Token{}, fmt.Errorf("token request failed with status %s", resp.Status)
	}

	var dto tokenDTO
	if err := json.NewDecoder(resp.Body).Decode(&dto); err != nil {
		return Token{}, fmt.Errorf("decode token response: %w", err)
	}

	return Token{
		AccessToken: dto.AccessToken,
		TokenType:   dto.TokenType,
		ExpiresAt:   time.Now().Add(time.Duration(dto.ExpiresIn) * time.Second),
	}, nil
}
