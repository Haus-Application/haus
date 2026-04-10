package nest

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Google OAuth2 endpoints for Nest/SDM.
// The auth URL goes through nestservices.google.com (not accounts.google.com!)
// because Device Access uses a partner connection flow. I learned this the hard
// way after spending three hours debugging why the standard Google OAuth URL
// wasn't returning SDM scopes. Three hours I'll never get back, Michael.
const (
	AuthURLTemplate = "https://nestservices.google.com/partnerconnections/%s/auth"
	TokenURL        = "https://www.googleapis.com/oauth2/v4/token"
)

// SDM API scope. This single scope covers all device read/write operations.
const SDMScope = "https://www.googleapis.com/auth/sdm.service"

// OAuthConfig holds the Google OAuth2 credentials and project settings.
// You get the ClientID and ClientSecret from the Google Cloud Console,
// and the ProjectID from the Device Access Console. They're different consoles.
// Google has a lot of consoles. It's like the banana stand -- there's always
// another one.
type OAuthConfig struct {
	ClientID     string
	ClientSecret string
	ProjectID    string
	RedirectURI  string
}

// TokenResponse is the response from the Google OAuth2 token endpoint.
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token,omitempty"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope,omitempty"`
}

// ExpiresAt returns the time when the access token will expire, calculated
// from the current time and the expires_in duration.
func (t *TokenResponse) ExpiresAt() time.Time {
	return time.Now().Add(time.Duration(t.ExpiresIn) * time.Second)
}

// GetAuthURL builds the Google OAuth2 authorization URL for the Nest SDM flow.
// The user opens this URL in a browser, authorizes the app, and gets redirected
// back with an authorization code. It's a whole thing. The redirect_uri must
// match exactly what's configured in the Google Cloud Console, down to the
// trailing slash. Ask me how I know.
func GetAuthURL(config OAuthConfig) string {
	authBase := fmt.Sprintf(AuthURLTemplate, config.ProjectID)

	params := url.Values{
		"redirect_uri":  {config.RedirectURI},
		"access_type":   {"offline"}, // Required to get a refresh_token
		"prompt":        {"consent"}, // Force consent screen to ensure refresh_token is returned
		"client_id":     {config.ClientID},
		"response_type": {"code"},
		"scope":         {SDMScope},
	}

	return authBase + "?" + params.Encode()
}

// ExchangeCode exchanges an authorization code for access and refresh tokens.
// This is the second step of the OAuth2 flow -- the user has authorized us,
// Google gave us a code, and now we trade that code for the real tokens.
// The code is single-use and expires after a few minutes, so don't dawdle.
func ExchangeCode(config OAuthConfig, code string) (*TokenResponse, error) {
	form := url.Values{
		"client_id":     {config.ClientID},
		"client_secret": {config.ClientSecret},
		"code":          {code},
		"grant_type":    {"authorization_code"},
		"redirect_uri":  {config.RedirectURI},
	}

	resp, err := http.Post(TokenURL, "application/x-www-form-urlencoded", strings.NewReader(form.Encode()))
	if err != nil {
		return nil, fmt.Errorf("nest: token exchange request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("nest: read token response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("nest: token exchange failed (status %d): %s", resp.StatusCode, body)
	}

	var token TokenResponse
	if err := json.Unmarshal(body, &token); err != nil {
		return nil, fmt.Errorf("nest: unmarshal token response: %w", err)
	}

	return &token, nil
}

// RefreshToken uses a refresh token to obtain a new access token. Access tokens
// expire after 3600 seconds (one hour), so you need to refresh before they expire.
// The refresh token itself doesn't expire unless the user revokes access, which
// is nice because it means we don't have to make them go through the whole
// authorization flow again. Mother would NOT appreciate being asked to re-authorize.
func RefreshToken(config OAuthConfig, refreshToken string) (*TokenResponse, error) {
	form := url.Values{
		"client_id":     {config.ClientID},
		"client_secret": {config.ClientSecret},
		"refresh_token": {refreshToken},
		"grant_type":    {"refresh_token"},
	}

	resp, err := http.Post(TokenURL, "application/x-www-form-urlencoded", strings.NewReader(form.Encode()))
	if err != nil {
		return nil, fmt.Errorf("nest: token refresh request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("nest: read refresh response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("nest: token refresh failed (status %d): %s", resp.StatusCode, body)
	}

	var token TokenResponse
	if err := json.Unmarshal(body, &token); err != nil {
		return nil, fmt.Errorf("nest: unmarshal refresh response: %w", err)
	}

	// Refresh responses don't include the refresh_token, so preserve the original.
	if token.RefreshToken == "" {
		token.RefreshToken = refreshToken
	}

	return &token, nil
}
