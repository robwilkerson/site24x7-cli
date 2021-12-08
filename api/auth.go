// Documentation of Site24x7 authentication patterns and process is located at
// https://www.site24x7.com/help/api/#authentication
// - Register an application and create a scoped grant token at https://api-console.zoho.com
// - Exchange the grant token for a refresh token

package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/spf13/viper"
)

// Token contains the data returned from a call to exchange either a grant or
// a refresh token for an access token.
type Token struct {
	AccessToken  string  `json:"access_token"`
	RefreshToken string  `json:"refresh_token"`
	ExpiresIn    int     `json:"expires_in"`
	ApiDomain    string  `json:"api_domain"`
	TokenType    string  `json:"token_type"`
	Error        *string `json:"error"`
}

// Configure accepts a Site24x7 grant token and returns a long-lived refresh
// token.
func Configure(grantToken string) (string, error) {
	t, err := exchangeToken(grantToken, "authorization_code")
	if err != nil {
		return "", err
	}

	return t, nil
}

// Authenticate fetches and stores a short-lived access token that will be used
// in subsequent API calls.
func Authenticate() {
	t, err := exchangeToken(viper.GetString("auth.refresh_token"), "refresh_token")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	os.Setenv("AUTH_ACCESS_TOKEN", t)
}

// exchangeToken exchanges a grant token (aka "authorization code") for a
// refresh token or a refresh token for an access token.
func exchangeToken(token string, grantType string) (string, error) {

	// Build the request

	endpoint := fmt.Sprintf("%s/oauth/v2/token", os.Getenv("AUTH_BASE_URL"))
	data := url.Values{
		"client_id":     {os.Getenv("AUTH_CLIENT_ID")},
		"client_secret": {os.Getenv("AUTH_CLIENT_SECRET")},
		"grant_type":    {grantType},
	}
	// Include the appropriate token
	switch grantType {
	case "authorization_code":
		data.Set("code", token)
	case "refresh_token":
		data.Set("refresh_token", token)
	}
	payload := strings.NewReader(data.Encode())
	req, err := http.NewRequest("POST", endpoint, payload)
	if err != nil {
		return "", fmt.Errorf("[Auth.exchangeToken] ERROR: Unable to create request (%s)", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Send the request

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("[Auth.exchangeToken] ERROR: unable to execute request (%s)", err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("[Auth.exchangeToken] ERROR: Unable to read response body (%s)", err)
	}

	// Unmarshal the response
	var t Token
	if err := json.Unmarshal(body, &t); err != nil {
		return "", fmt.Errorf("[Auth.exchangeToken] ERROR: Unable to  parse response body (%s)", err)
	}
	if t.Error != nil {
		return "", fmt.Errorf("[Auth.exchangeToken] ERROR: received an error response from Site24x7 (%s)", *t.Error)
	}

	// Return the appropriate exchange token
	switch grantType {
	case "authorization_code":
		return t.RefreshToken, nil
	case "refresh_token":
		return t.AccessToken, nil
	default:
		return "", fmt.Errorf("[Auth.exchangeToken] ERROR: Unrecognized grant type (%s)", grantType)
	}
}

// httpHeader returns the header name and value required to authenticate each
// API request sent to Site24x7.
func httpHeader() (string, string) {
	return "Authorization", fmt.Sprintf("Zoho-oauthtoken %s", os.Getenv("AUTH_ACCESS_TOKEN"))
}
