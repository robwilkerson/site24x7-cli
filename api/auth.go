// Documentation of Site24x7 authentication patterns and process is located at
// https://www.site24x7.com/help/api/#authentication
// - Register an application and create a scoped grant token at https://api-console.zoho.com
// - Exchange the grant token for a refresh token

package api

import (
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/spf13/viper"
)

// AuthToken contains the data returned from a call to exchange either a grant or
// a refresh token for an access token.
type AuthToken struct {
	AccessToken  string  `json:"access_token"`
	RefreshToken string  `json:"refresh_token"`
	ExpiresIn    int     `json:"expires_in"`
	APIDomain    string  `json:"api_domain"`
	TokenType    string  `json:"token_type"`
	Error        *string `json:"error"`
}

// Configure exchanges a short-lived grant token (a.k.a. authorization code) and
// returns a long-lived refresh token.
func Configure(grantToken string) (string, error) {
	exchangableToken := map[string]string{
		"grantType": "authorization_code",
		"key":       "code",
		"value":     grantToken,
	}

	t, err := exchangeToken(exchangableToken)
	if err != nil {
		return "", err
	}

	return t.RefreshToken, nil
}

// Authenticate exchanges a refresh token for a short-lived access token and
// stores the latter for use in subsequent API calls.
func Authenticate() error {
	exchangableToken := map[string]string{
		"grantType": "refresh_token",
		"key":       "refresh_token",
		"value":     viper.GetString("auth.refresh_token"),
	}

	t, err := exchangeToken(exchangableToken)
	if err != nil {
		return err
	}

	os.Setenv("AUTH_ACCESS_TOKEN", t.AccessToken)

	return nil
}

// exchangeToken exchanges a grant token (aka "authorization code") for a
// refresh token or a refresh token for an access token.
func exchangeToken(token map[string]string) (*AuthToken, error) {
	req := Request{
		Endpoint: fmt.Sprintf("%s/oauth/v2/token", os.Getenv("AUTH_BASE_URL")),
		Method:   "POST",
		Headers: http.Header{
			"Content-Type": {"application/x-www-form-urlencoded"},
		},
		Body: nil,
		QueryString: url.Values{
			"client_id":     {viper.GetString("auth.client_id")},
			"client_secret": {viper.GetString("auth.client_secret")},
			"grant_type":    {token["grantType"]},
			token["key"]:    {token["value"]},
		},
	}

	t, err := req.FetchAuthToken()
	if err != nil {
		return nil, err
	}
	if t.Error != nil {
		return nil, fmt.Errorf("[Auth.exchangeToken] ERROR: received an error response from Site24x7 (%s)", *t.Error)
	}

	return t, nil
}

// httpHeader returns the header name and value required to authenticate each
// API request sent to Site24x7.
func httpHeader() (string, string) {
	return "Authorization", fmt.Sprintf("Zoho-oauthtoken %s", os.Getenv("AUTH_ACCESS_TOKEN"))
}
