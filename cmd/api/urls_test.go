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
	if os.Getenv("INTEGRATION") == "" {
		t.Skip("skipping integration tests")
	}

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
	if os.Getenv("INTEGRATION") == "" {
		t.Skip("skipping integration tests")
	}

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
	if os.Getenv("INTEGRATION") == "" {
		t.Skip("skipping integration tests")
	}

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
	if os.Getenv("INTEGRATION") == "" {
		t.Skip("skipping integration tests")
	}

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
