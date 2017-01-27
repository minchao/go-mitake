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
	LibraryVersion   = "0.0.1"
	DefaultBaseURL   = "https://smexpress.mitake.com.tw:9601"
	DefaultUserAgent = "go-mitake/" + LibraryVersion
)

// NewClient returns a new Mitake API client.
func NewClient(username, password string, httpClient *http.Client) *Client {
	if username == "" || password == "" {
		log.Fatal("username or password cannot be empty")
	}
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	baseURL, _ := url.Parse(DefaultBaseURL)

	return &Client{
		client:    httpClient,
		username:  username,
		password:  password,
		UserAgent: DefaultUserAgent,
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
	if c := r.StatusCode; 200 <= c && c <= 299 {
		if r.ContentLength == 0 {
			return errors.New("unexpected empty body")
		}
		return nil
	} else {
		// Mitake API always return status code 200
		return fmt.Errorf("unexpected status code: %d", c)
	}
}

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

func (c *Client) Get(url string) (*http.Response, error) {
	req, err := c.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}

func (c *Client) Post(url string, bodyType string, body io.Reader) (*http.Response, error) {
	req, err := c.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", bodyType)
	return c.Do(req)
}

func (c *Client) buildDefaultQuey() url.Values {
	q := url.Values{}
	q.Set("username", c.username)
	q.Set("password", c.password)
	return q
}
