// Documentation of Site24x7 authentication patterns and process is located at
// https://www.site24x7.com/help/api/#authentication

package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

// AccessToken defines a short-lived token that can be fetched, then used in
// API calls.
type AccessToken struct {
	AccessToken string  `json:"access_token"`
	ExpiresIn   int     `json:"expires_in"`
	ApiDomain   string  `json:"api_domain"`
	TokenType   string  `json:"token_type"`
	Error       *string `json:"error"`
}

// Authenticate fetches and stores a short-lived access token that will be used
// in subsequent API calls.
func Authenticate() {
	at, err := generateAccessToken(os.Getenv("AUTH_REFRESH_TOKEN"))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	os.Setenv("AUTH_ACCESS_TOKEN", at)
}

// generateAccessToken exchanges a long-lived refresh token that must be created
// manually for a short-lived access token.
func generateAccessToken(refreshToken string) (string, error) {

	// Build the request

	endpoint := fmt.Sprintf("%s/oauth/v2/token", os.Getenv("AUTH_BASE_URL"))
	data := url.Values{
		"client_id":     {os.Getenv("AUTH_CLIENT_ID")},
		"client_secret": {os.Getenv("AUTH_CLIENT_SECRET")},
		"refresh_token": {refreshToken},
		"grant_type":    {"refresh_token"},
	}
	payload := strings.NewReader(data.Encode())
	req, err := http.NewRequest("POST", endpoint, payload)
	if err != nil {
		return "", fmt.Errorf("[generateAccessToken] ERROR: Unable to create request (%s)", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Send the request

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("[generateAccessToken] ERROR: unable to execute request (%s)", err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("[generateAccessToken] ERROR: Unable to read response body (%s)", err)
	}

	// Unmarshal the response

	var at AccessToken
	if err := json.Unmarshal(body, &at); err != nil {
		return "", fmt.Errorf("[generateAccessToken] ERROR: Unable to  parse response body (%s)", err)
	}
	if at.Error != nil {
		return "", fmt.Errorf("[generateAccessToken] ERROR: received an error response from Site24x7 (%s)", *at.Error)
	}

	return at.AccessToken, nil
}

// httpHeader returns the header name and value required to authenticate each
// API request sent to Site24x7.
func httpHeader() (string, string) {
	return "Authorization", fmt.Sprintf("Zoho-oauthtoken %s", os.Getenv("AUTH_ACCESS_TOKEN"))
}
