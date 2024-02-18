package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"site24x7/logger"
	"strings"
)

// Request provides any request-specific data that might be required to access
// a Site24x7 endpoint.
type Request struct {
	Endpoint    string
	Method      string
	Body        []byte
	Headers     http.Header
	QueryString url.Values
}

// ApiResponse defines the top level schema of (almost?) every Site24x7 API
// response. The Data component contains the domain model data and varies wildly
// between API calls, so it must be flexible. We'll handle it as raw JSON at
// this stage and unmarshal it separately when we need it.
type ApiResponse struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

// EmptyApiResponse defines the schema of a response that includes no data. One
// example of such a response comes from the POST /milestone endpoint.
// https://www.site24x7.com/help/api/#add-a-milestone-marker
type EmptyApiResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// FetchToken returns a refresh token, an access token, or both depending on
// whether a grant token or a refresh token is being exchanged.
//
// This endpoint is also a bit of an odd duck in that data is posted, but the
// body contains serialized form data rather than json. The API response itself
// also doesn't return the traditional schema so it just feels a little cleaner
// to handle it separately from other fetches.
//
// This function is called before any command is executed and cannot reliably
// use the logger.
func (r *Request) FetchAuthToken() (*AuthToken, error) {
	// Weirdness #1: serialize the query string data so it can be sent as the
	// request body. ¯\_(ツ)_/¯
	body := strings.NewReader(r.QueryString.Encode())
	req, err := http.NewRequest(r.Method, r.Endpoint, body)
	if err != nil {
		return nil, fmt.Errorf("[api.FetchAuthToken] ERROR: Unable to create request (%s)", err)
	}
	req.Header = r.Headers

	// Uncomment to debug requests
	// dump, _ := httputil.DumpRequestOut(req, true)
	// fmt.Printf("%q\n", dump)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("[api.FetchAuthToken] ERROR: unable to execute request (%s)", err)
	}
	defer res.Body.Close()

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("[api.FetchAuthToken] ERROR: Unable to read response body (%s)", err)
	}

	// Weirdness #2: the endpoint returns a token, but not in ApiResponse.Data
	// like, well, every other endpoint I've tried thus far.
	var t AuthToken
	if err := json.Unmarshal(b, &t); err != nil {
		return nil, fmt.Errorf("[api.FetchAuthToken] ERROR: Unable to  parse response body (%s)", err)
	}

	return &t, nil
}

// Fetch calls a Site24x7 API and returns the response.
func (r *Request) Fetch() (*ApiResponse, error) {
	body := bytes.NewReader(r.Body)
	qs := strings.NewReader(r.QueryString.Encode())

	var req *http.Request
	if r.Method == "GET" {
		req, _ = http.NewRequest(r.Method, r.Endpoint, qs)
	} else {
		req, _ = http.NewRequest(r.Method, r.Endpoint, body)
	}
	req.Header = r.Headers

	// Apply common headers
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")

	// TODO: only do the work if verbosity == debug?
	dumpreq, _ := httputil.DumpRequestOut(req, true)
	logger.Debug(fmt.Sprintf("[api.Fetch] Request: %q\n", dumpreq))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("[api.Fetch] ERROR: unable to execute request (%s)", err)
	}
	defer res.Body.Close()

	// TODO: only do the work if verbosity == debug?
	dumpres, _ := httputil.DumpResponse(res, true)
	logger.Debug(fmt.Sprintf("[api.Fetch] Response: %q\n", dumpres))

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("[api.Fetch] ERROR: Unable to read response body (%s)", err)
	}

	var ar ApiResponse
	if err := json.Unmarshal(b, &ar); err != nil {
		return nil, fmt.Errorf("[api.Fetch] ERROR: Unable to  parse response body (%s)", err)
	}

	return &ar, nil
}
