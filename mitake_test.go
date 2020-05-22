package mitake

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"testing"
)

func setup() (client *Client, mux *http.ServeMux, teardown func()) {
	// mux is the HTTP request multiplexer used with the test server.
	mux = http.NewServeMux()

	// server is a test HTTP server used to provide mock API responses.
	server := httptest.NewServer(mux)

	// mitake client configured to use test server
	baseURL, _ := url.Parse(server.URL)

	// client is the mitake client being tested.
	client = NewClient("username", "password", nil)
	client.BaseURL = baseURL

	return client, mux, server.Close
}

func testMethod(t *testing.T, r *http.Request, want string) {
	if got := r.Method; got != want {
		t.Errorf("Request method is %v, want %v", got, want)
	}
}

func testINI(t *testing.T, r *http.Request, want string) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		t.Errorf("Request parameters error: %v", err)
	}
	defer r.Body.Close()

	if got := string(body); got != want {
		t.Errorf("Request parameters is %v, want %v", got, want)
	}
}

func TestNewClient(t *testing.T) {
	c := NewClient("username", "password", nil)

	if got, want := c.BaseURL.String(), defaultBaseURL; got != want {
		t.Errorf("NewClient BaseURL is %v, want %v", got, want)
	}
}

func TestClient_NewRequest(t *testing.T) {
	c := NewClient("username", "password", nil)

	inURL, outURL := "/foo", defaultBaseURL+"foo"
	inBody, outBody := "Hello, 世界", "Hello, 世界"
	req, _ := c.NewRequest("GET", inURL, strings.NewReader(inBody))

	if got, want := req.URL.String(), outURL; got != want {
		t.Errorf("NewRequest(%q) URL is %v, want %v", inURL, got, want)
	}

	body, _ := ioutil.ReadAll(req.Body)
	if got, want := string(body), outBody; got != want {
		t.Errorf("NewRequest(%q) Body is %v, want %v", inBody, got, want)
	}

	if got, want := req.Header.Get("User-Agent"), c.UserAgent; got != want {
		t.Errorf("NewRequest() User-Agent is %v, want %v", got, want)
	}
}

func TestClient_Do(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if m := "GET"; m != r.Method {
			t.Errorf("Request method is %v, want %v", r.Method, m)
		}
		_, _ = fmt.Fprint(w, "Hello, 世界")
	})

	req, _ := client.NewRequest("GET", "/", nil)
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("Do returned unexpected error: %v", err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	want := "Hello, 世界"
	if !reflect.DeepEqual(string(body), want) {
		t.Errorf("Response body is %s, want %s", body, want)
	}
}

func TestClient_Do_httpError(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Bad Request", 400)
	})

	req, _ := client.NewRequest("GET", "/", nil)
	_, err := client.Do(req)

	if err == nil {
		t.Error("Expected HTTP 400 error.")
	}
}

func TestClient_Do_noContent(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, "")
	})

	req, _ := client.NewRequest("GET", "/", nil)
	_, err := client.Do(req)

	if err == nil {
		t.Error("Expected empty body error.")
	}
}
