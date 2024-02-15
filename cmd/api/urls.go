package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/asirago/shorturl/internal/database"
	"github.com/google/uuid"
)

type request struct {
	URL            string `json:"url"`
	CustomShortURL string `json:"custom_short_url"`
}

type response struct {
	URL            string        `json:"url"`
	CustomShortURL string        `json:"custom_short_url"`
	Expiry         time.Duration `json:"expiry"`
}

func (app *application) shortenUrl(w http.ResponseWriter, r *http.Request) {

	input := request{}

	rdb := database.CreateRedisClient()

	app.readJSON(w, r, &input)

	if input.CustomShortURL == "" {
		input.CustomShortURL = uuid.New().String()
	}
	// TODO: Check API ratelimit

	// TODO: Check if CustomShortURL is already in use

	err := rdb.Set(context.Background(), input.CustomShortURL, input.URL, 0).Err()
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

	resp := response{
		URL:            input.URL,
		CustomShortURL: input.CustomShortURL,
	}

	err = app.writeJSON(w, r, http.StatusCreated, resp, nil)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

}
