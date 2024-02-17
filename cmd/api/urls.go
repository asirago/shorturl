package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/asirago/shorturl/internal/database"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type request struct {
	URL            string `json:"url"`
	CustomShortURL string `json:"custom_short_url"`
}

type response struct {
	URL            string        `json:"url"`
	CustomShortURL string        `json:"custom_short_url"`
	Expiry         time.Duration `json:"expiry"`
	RateLimit      int           `json:"rate_limit"`
}

func (app *application) shortenUrl(w http.ResponseWriter, r *http.Request) {

	input := request{}
	ctx := context.Background()
	url_quota := 5
	rateLimitDuration := 24 * time.Hour

	rdb := database.CreateRedisClient()

	err := app.readJSON(w, r, &input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rateLimit, err := rdb.Get(ctx, host).Int()
	if err == redis.Nil {
		rdb.Set(ctx, host, url_quota-1, rateLimitDuration).Err()
		rateLimit = url_quota - 1
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		rateLimit64, err := rdb.Decr(ctx, host).Result()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		rateLimit = int(rateLimit64)
	}

	if rateLimit <= 0 {
		ttl, err := rdb.TTL(ctx, host).Result()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Error(
			w,
			fmt.Sprintf(
				"Rate limit exceeded, try again after %s\n",
				time.Now().UTC().Add(ttl).Format("2006-01-02 15:04:05 UTC"),
			),
			http.StatusServiceUnavailable,
		)
		return
	}

	// TODO: Check if CustomShortURL is already in use
	if input.CustomShortURL == "" {
		input.CustomShortURL = uuid.New().String()[:6]
	}

	err = rdb.Set(context.Background(), input.CustomShortURL, input.URL, 0).Err()
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

	resp := response{
		URL:            input.URL,
		CustomShortURL: input.CustomShortURL,
		RateLimit:      rateLimit,
	}

	err = app.writeJSON(w, r, http.StatusCreated, resp, nil)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

}
