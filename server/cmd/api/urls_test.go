package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/go-chi/chi/v5"
)

func TestShortenUrl(t *testing.T) {
	integrationTest(t)

	var request struct {
		URL            string `json:"url"`
		CustomShortURL string `json:"custom_short_url"`
	}
	request.URL = "https://google.com"
	request.CustomShortURL = ""

	b, _ := json.Marshal(request)

	r, _ := http.NewRequest(http.MethodPost, "/shorten-url", bytes.NewBuffer(b))
	r.RemoteAddr = "192.168.65.1:58298"

	rr := httptest.NewRecorder()

	var s *Server

	s.shortenUrl(rr, r)

	if rr.Code != http.StatusCreated {
		t.Fatalf("got %v, wanted %v", rr.Code, http.StatusCreated)
	}
}

func TestShortenUrlSameCustomURL(t *testing.T) {
	integrationTest(t)

	var request struct {
		URL            string `json:"url"`
		CustomShortURL string `json:"custom_short_url"`
	}
	request.URL = "https://google.com"
	request.CustomShortURL = "helloWorld"

	for i := 0; i < 2; i++ {
		b, _ := json.Marshal(request)

		r, _ := http.NewRequest(http.MethodPost, "/shorten-url", bytes.NewBuffer(b))
		r.RemoteAddr = "192.168.65.1:58298"

		rr := httptest.NewRecorder()

		var app *Server

		app.shortenUrl(rr, r)

		res := rr.Body.String()
		if i == 1 {
			if rr.Code != http.StatusConflict {
				t.Fatalf("got %v, want %v", rr.Code, http.StatusConflict)
			}

			if res != "Short custom url already taken\n" {
				t.Fatalf("got %v, want %v", res, "Short custom url already taken\n")
			}
		}
	}

}

func TestShortenUrlRateLimit(t *testing.T) {
	integrationTest(t)

	var request struct {
		URL            string `json:"url"`
		CustomShortURL string `json:"custom_short_url"`
	}
	request.URL = "https://google.com"
	request.CustomShortURL = ""

	for i := 0; i < 6; i++ {
		b, _ := json.Marshal(request)

		r, _ := http.NewRequest(http.MethodPost, "/shorten-url", bytes.NewBuffer(b))
		r.RemoteAddr = "127.0.0.1"
		rr := httptest.NewRecorder()

		var app *Server

		app.shortenUrl(rr, r)

		if i == 5 {
			if rr.Code != http.StatusServiceUnavailable {
				t.Fatalf("got %d, want %d\n", rr.Code, http.StatusServiceUnavailable)
			}
		}
	}
}

func TestResolveUrlThatDoesNotExist(t *testing.T) {
	integrationTest(t)

	testUrl := "/abcd"
	wantErr := fmt.Sprintf("%s does not exist", testUrl)

	var s *Server
	router := chi.NewRouter()
	router.Get("/{url}", s.resolveUrl)

	r, _ := http.NewRequest(http.MethodGet, testUrl, nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, r)

	if rr.Code != http.StatusNotFound {
		t.Errorf("want %d, but got %d", http.StatusNotFound, rr.Code)
	}

	var resp struct {
		Error string `json:"error"`
	}

	js := json.NewDecoder(rr.Body)
	err := js.Decode(&resp)
	if err != nil {
		t.Fatal("could not decode json")
	}

	if resp.Error != wantErr {
		t.Errorf(wantErr)
	}
}

func TestResolveUrlExists(t *testing.T) {
	integrationTest(t)

	// set up chi router
	var s *Server
	router := chi.NewRouter()
	router.Post("/shorten-url", s.shortenUrl)
	router.Get("/{url}", s.resolveUrl)

	// Shorten url
	var request struct {
		URL            string `json:"url"`
		CustomShortURL string `json:"custom_short_url"`
	}
	request.URL = "https://google.se"

	b, _ := json.Marshal(request)

	r, _ := http.NewRequest(http.MethodPost, "/shorten-url", bytes.NewBuffer(b))
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, r)

	var resp struct {
		URL      string `json:"url"`
		ShortURL string `json:"short_url"`
	}
	js := json.NewDecoder(rr.Body)
	err := js.Decode(&resp)
	if err != nil {
		t.Fatal("could not decode json response")
	}

	r, _ = http.NewRequest(http.MethodGet, fmt.Sprintf("/%s", resp.ShortURL), nil)
	rr = httptest.NewRecorder()

	router.ServeHTTP(rr, r)

	var resp2 struct {
		URL string `json:"url"`
	}

	err = json.NewDecoder(rr.Body).Decode(&resp2)
	if err != nil {
		t.Fatal("could not decore json response")
	}

	if resp.URL != resp2.URL {
		t.Errorf("want %s, but got %s", resp.URL, resp2.URL)
	}
}

func TestCleanUrl(t *testing.T) {
	testCases := []struct {
		name    string
		testUrl string
		wantUrl string
		wantErr string
	}{
		{
			"Url",
			"https://google.se",
			"https://google.se",
			"",
		},
		{
			"UrlWithSubdomain",
			"https://www.tinkaling.asirago.xyz",
			"https://www.tinkaling.asirago.xyz",
			"",
		},
		{
			"UrlWithPaths",
			"https://facebook.com/foo/bar",
			"https://facebook.com/foo/bar",
			"",
		},
		{
			"UrlWithDoubleExtensionTopLevelDomain",
			"https://amazon.co.uk",
			"https://amazon.co.uk",
			"",
		},
		{
			"UrlWithHttpScheme",
			"https://amazon.com",
			"https://amazon.com",
			"http scheme not allowed, only supports https scheme",
		},
		{
			"UrlWithQueryParameters",
			"amazon.co.uk/foo/bar?color=blue&colour=green",
			"https://amazon.co.uk/foo/bar?color%3Dblue%26colour%3Dgreen",
			"",
		},
		{
			"UrlWithTrailingSlash",
			"https://example.com/foo/bar/",
			"https://example.com/foo/bar",
			"<script> forbidden",
		},
		{
			"UrlWithScriptTag",
			"example.com/<script>alert('XSS')</script>",
			"",
			"<script> forbidden",
		},
		{
			"UrlWithFragment",
			"example.com/foo/bar?color=green#hello",
			"https://example.com/foo/bar?color%3Dgreen#hello",
			"",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gotUrl, err := cleanURL(tc.testUrl)
			if err != nil && err.Error() != tc.wantErr {
				t.Errorf("want %s, but got %s", tc.wantErr, err.Error())
			}

			if tc.wantUrl != gotUrl {
				t.Errorf("want %s, but got %s", tc.wantUrl, gotUrl)
			}

		})
	}
}

func integrationTest(t *testing.T) {
	t.Helper()
	if os.Getenv("INTEGRATION") == "" {
		t.Skip("skipping integration tests: set environment variable INTEGRATION")
	}

}
