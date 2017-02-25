package unit

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"testing"

	"github.com/chibisov/go-yadisk/yadisk"
)

func TestNewClient(t *testing.T) {
	c := yadisk.NewClient("ACCESS_TOKEN")

	if got, want := c.HTTPClient, http.DefaultClient; got != want {
		t.Errorf("NewClient HTTPClient is %v, want %v", got, want)
	}

	if got, want := c.BaseURL.String(), "https://cloud-api.yandex.net/"; got != want {
		t.Errorf("NewClient BaseURL is %v, want %v", got, want)
	}

	if got, want := c.AccessToken, "ACCESS_TOKEN"; got != want {
		t.Errorf("NewClient AccessToken is %v, want %v", got, want)
	}
}

func TestNewRequest(t *testing.T) {
	c := yadisk.NewClient("ACCESS_TOKEN")

	type User struct {
		Login string `json:"login"`
	}

	inURL, outURL := "disk", "https://cloud-api.yandex.net/v1/disk/"
	inBody, outBody := &User{Login: "sosisa"}, `{"login":"sosisa"}`+"\n"
	req, _ := c.NewRequest("GET", inURL, inBody)

	// Check that relative URL was expanded.
	if got, want := req.URL.String(), outURL; got != want {
		t.Errorf("NewRequest(%q) URL is %v, want %v", inURL, got, want)
	}

	// Check that body was JSON encoded.
	body, _ := ioutil.ReadAll(req.Body)
	if got, want := string(body), outBody; got != want {
		t.Errorf("NewRequest(%q) Body is %v, want %v", inBody, got, want)
	}

	// Check that authorization key is provided.
	got, want := req.Header.Get("Authorization"), "OAuth ACCESS_TOKEN"
	if got != want {
		t.Errorf("Authorization header is %v, want %v", got, want)
	}
}

func TestNewRequest_invalid_json(t *testing.T) {
	c := yadisk.NewClient("ACCESS_TOKEN")

	type User struct {
		Data map[interface{}]interface{}
	}
	_, err := c.NewRequest("GET", "/", &User{})

	if err == nil {
		t.Error("Expected error to be returned")
	}
	if err, ok := err.(*json.UnsupportedTypeError); !ok {
		t.Errorf("Expected a JSON error; got %#v", err)
	}
}

func TestNewRequest_bad_url(t *testing.T) {
	c := yadisk.NewClient("ACCESS_TOKEN")
	_, err := c.NewRequest("GET", ":", nil)

	if err == nil {
		t.Errorf("Expected error to be returned")
	}
	if err, ok := err.(*url.Error); !ok {
		t.Errorf("Expected URL parse error, got %+v", err)
	}
}

// If a nil body is passed to github.NewRequest, make sure that nil is also
// passed to http.NewRequest. In most cases, passing an io.Reader that returns
// no content is fine, since there is no difference between an HTTP request
// body that is an empty string versus one that is not set at all. However in
// certain cases, intermediate systems may treat these differently resulting in
// subtle errors.
func TestNewRequest_empty_body(t *testing.T) {
	c := yadisk.NewClient("ACCESS_TOKEN")
	req, err := c.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatalf("NewRequest returned unexpected error: %v", err)
	}
	if req.Body != nil {
		t.Fatalf("constructed request contains a non-nil Body: %v", req.Body)
	}
}

func TestDo(t *testing.T) {
	setup()
	defer teardown()

	buf := &bytes.Buffer{}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if m := "GET"; m != r.Method {
			t.Errorf("Request method = %v, want %v", r.Method, m)
		}
		fmt.Fprint(w, `{"A":"a"}`)
	})

	req, _ := client.NewRequest("GET", "/", nil)
	client.Do(context.Background(), req, buf)

	if got, want := buf.String(), `{"A":"a"}`; got != want {
		t.Errorf("Response body = %v, want %v", got, want)
	}
}

func TestDo_with_struct(t *testing.T) {
	setup()
	defer teardown()

	type foo struct {
		A string
	}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if m := "GET"; m != r.Method {
			t.Errorf("Request method = %v, want %v", r.Method, m)
		}
		fmt.Fprint(w, `{"A":"a"}`)
	})

	req, _ := client.NewRequest("GET", "/", nil)
	body := new(foo)
	client.Do(context.Background(), req, body)

	want := &foo{A: "a"}
	if !reflect.DeepEqual(body, want) {
		t.Errorf("Response body = %v, want %v", body, want)
	}
}

func TestDo_http_error_not_json(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Bad Request", 400)
	})

	req, _ := client.NewRequest("GET", "/", nil)
	_, err := client.Do(context.Background(), req, nil)

	if err == nil {
		t.Error("Expected HTTP 400 error.")
	}

	if _, ok := err.(*yadisk.APIError); !ok {
		t.Errorf("Expected a yadisk.APIError error; got %#v", err)
	}

	if want, got := "Yandex.Disk API error", err.Error(); got != want {
		t.Errorf("Error text = %s, want %s", got, want)
	}
}

func TestDo_http_error_json(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Error(
			w,
			`
			{
				"description": "resource already exists",
				"error": "PlatformResourceAlreadyExists"
			}
			`,
			409,
		)
	})

	req, _ := client.NewRequest("GET", "/", nil)
	_, err := client.Do(context.Background(), req, nil)

	if err == nil {
		t.Error("Expected HTTP 409 error.")
	}

	if _, ok := err.(*yadisk.APIError); !ok {
		t.Errorf("Expected a yadisk.APIError error; got %#v", err)
	}

	want := "Yandex.Disk API error. Code: PlatformResourceAlreadyExists. " +
		"Description: resource already exists."
	got := err.Error()
	if got != want {
		t.Errorf("Error text = '%s', want '%s'", got, want)
	}
}
