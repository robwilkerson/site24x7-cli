package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// endpoint, method, custom headers
func (t *Token) Fetch(endpoint string, method string, body *strings.Reader, headers http.Header) error {
	req, err := http.NewRequest(method, endpoint, body)
	if err != nil {
		return fmt.Errorf("[Token.Fetch] ERROR: Unable to create request (%s)", err)
	}
	req.Header = headers

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("[Token.Fetch] ERROR: unable to execute request (%s)", err)
	}
	defer res.Body.Close()

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("[Token.Fetch] ERROR: Unable to read response body (%s)", err)
	}

	// Unmarshal the response
	if err := json.Unmarshal(b, t); err != nil {
		return fmt.Errorf("[Auth.exchangeToken] ERROR: Unable to  parse response body (%s)", err)
	}

	return nil
}
