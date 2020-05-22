package mitake

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

const (
	libraryVersion   = "2.08.1"
	defaultUserAgent = "go-mitake/" + libraryVersion
	defaultBaseURL   = "https://smsapi.mitake.com.tw:443/"
)

// NewClient returns a new Mitake API client. The username and password are required
// for authentication. If a nil httpClient is provided, http.DefaultClient will be used.
func NewClient(username, password string, httpClient *http.Client) *Client {
	if username == "" || password == "" {
		log.Fatal("username or password cannot be empty")
	}
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	baseURL, _ := url.Parse(defaultBaseURL)

	return &Client{
		client:    httpClient,
		username:  username,
		password:  password,
		UserAgent: defaultUserAgent,
		BaseURL:   baseURL,
	}
}

// A Client manages communication with the Mitake API.
type Client struct {
	client   *http.Client
	username string
	password string

	BaseURL   *url.URL
	UserAgent string
}

// checkErrorResponse checks the API response for errors.
func checkErrorResponse(r *http.Response) error {
	c := r.StatusCode
	if 200 <= c && c <= 299 {
		if r.ContentLength == 0 {
			return errors.New("unexpected empty body")
		}
		return nil
	}
	// Mitake API always return status code 200
	return fmt.Errorf("unexpected status code: %d", c)
}

// Do sends an API request, and returns the API response.
// If the returned error is nil, the Response will contain a non-nil
// Body which the user is expected to close.
func (c *Client) Do(req *http.Request) (*http.Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	if err := checkErrorResponse(resp); err != nil {
		resp.Body.Close()
		return resp, err
	}
	return resp, nil
}

// NewRequest creates an API request. A relative URL can be provided in urlStr,
// in which case it is resolved relative to the BaseURL of the Client.
// Relative URLs should always be specified without a preceding slash.
func (c *Client) NewRequest(method, urlStr string, body io.Reader) (*http.Request, error) {
	rel, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	u := c.BaseURL.ResolveReference(rel)

	req, err := http.NewRequest(method, u.String(), body)
	if err != nil {
		return nil, err
	}

	if c.UserAgent != "" {
		req.Header.Set("User-Agent", c.UserAgent)
	}
	return req, nil
}

// Get method make a GET HTTP request.
func (c *Client) Get(url string) (*http.Response, error) {
	req, err := c.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}

// Post method make a POST HTTP request.
func (c *Client) Post(url string, bodyType string, body io.Reader) (*http.Response, error) {
	req, err := c.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", bodyType)
	return c.Do(req)
}

// buildDefaultQuery returns the default query string with authentication parameters.
func (c *Client) buildDefaultQuery() url.Values {
	q := url.Values{}
	q.Set("username", c.username)
	q.Set("password", c.password)
	return q
}
