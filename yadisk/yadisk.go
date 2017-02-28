package yadisk

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	"golang.org/x/net/context/ctxhttp"
)

const (
	defaultBaseURL = "https://cloud-api.yandex.net/"
	apiVersion     = "1"
)

// A Client manages communication with the Yandex.Disk API.
type Client struct {
	HTTPClient *http.Client // HTTP client used to communicate with the API.

	AccessToken string // OAuth access token

	// Base URL for API requests. Defaults to the public Yandex.Disk API.
	BaseURL *url.URL

	// Services used for talking to different parts of the Yandex.Disk API.
	Disk      *DiskService
	Resources *ResourcesService
}

type service struct {
	client *Client
}

// NewClient returns a new Yandex.Disk API client. If a nil httpClient is
// provided, http.DefaultClient will be used. To use API methods which require
// authentication, provide an http.Client that will perform the authentication
// for you (such as that provided by the golang.org/x/oauth2 library).
func NewClient(accessToken string) *Client {
	baseURL, _ := url.Parse(defaultBaseURL)

	c := &Client{
		HTTPClient:  http.DefaultClient,
		BaseURL:     baseURL,
		AccessToken: accessToken,
	}
	c.Disk = &DiskService{client: c}
	c.Resources = &ResourcesService{client: c}

	return c
}

// NewRequest creates an API request. A relative URL can be provided
// in urlStr, which will be resolved to the BaseURL of the Client.
// Relative URLS should always be specified without a preceding slash.
// If specified, the value pointed to by body is JSON encoded
// and included in as the request body.
func (c *Client) NewRequest(method, urlStr string, body interface{}) (*http.Request, error) {
	rel, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	// Build the request url.
	rel.Path = "v" + apiVersion + "/" + rel.Path + "/"
	u := c.BaseURL.ResolveReference(rel)

	// Build the request body.
	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err = json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	// Build the http request.
	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	// Set the necessary headers.
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "OAuth "+c.AccessToken)

	return req, nil
}

// Do sends an API request and returns the API response.
// The API response is JSON decoded and stored in the value
// pointed to by v, or returned as an error if an API error has occurred.
// If v implements the io.Writer interface,
// the raw response will be written to v, without attempting to decode it.
// The provided ctx must be non-nil. If it is canceled or times out,
// ctx.Err() will be returned.
func (c *Client) Do(ctx context.Context, req *http.Request, v interface{}) (*http.Response, error) {
	// Make the http request.
	resp, err := ctxhttp.Do(ctx, c.HTTPClient, req) // todo: test context
	if err != nil {
		return nil, err
	}

	// Defer the body closing.
	defer func() {
		// Drain up to 512 bytes and close the body to
		// let the Transport reuse the connection.
		// Read more https://groups.google.com/forum/#!topic/golang-nuts/4Rr8BYVKrAI
		io.CopyN(ioutil.Discard, resp.Body, 512)
		resp.Body.Close()
	}()

	// Check for the response errors.
	if err = checkResponse(resp); err != nil {
		return resp, err
	}

	// Fill the v variable with the response data if it's provided.
	if v != nil {
		if w, ok := v.(io.Writer); ok {
			// Write to the buffer if io.Writer is provided.
			_, err = io.Copy(w, resp.Body)
			if err != nil {
				return nil, err
			}
		} else {
			// Decode JSON to the struct if struct is provided.
			err = json.NewDecoder(resp.Body).Decode(v)
			if err == io.EOF {
				err = nil // ignore EOF errors caused by empty response body
			}
		}
	}

	return resp, err
}

// checkResponse checks the API response for errors,
// and returns them if present.
func checkResponse(r *http.Response) error {
	if r.StatusCode >= 400 {
		apiErr := new(APIError)

		// Skipping the json decoding errors.
		json.NewDecoder(r.Body).Decode(apiErr)

		return apiErr
	}

	return nil
}
