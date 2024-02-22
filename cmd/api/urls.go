package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/asirago/shorturl/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type request struct {
	URL            string `json:"url"`
	CustomShortURL string `json:"custom_short_url"`
}

type response struct {
	URL            string `json:"url"`
	CustomShortURL string `json:"custom_short_url"`
	RateLimit      int    `json:"rate_limit"`
}

func (s *Server) resolveUrl(w http.ResponseWriter, r *http.Request) {
	url := chi.URLParam(r, "url")

	rdb := database.CreateRedisClient()

	longUrl, err := rdb.Get(database.Ctx, url).Result()
	if err == redis.Nil {
		s.errorJSON(w, http.StatusNotFound, fmt.Sprintf("/%s does not exist", url))
		return
	} else if err != nil {
		s.serverErrorResponse(w, err)
		return
	}

	err = s.writeJSON(w, http.StatusMovedPermanently, map[string]string{"url": longUrl}, nil)
	if err != nil {
		s.serverErrorResponse(w, err)
	}
}

// TODO: REFACTOR
func (s *Server) shortenUrl(w http.ResponseWriter, r *http.Request) {

	req := request{}
	resp := response{}
	ctx := database.Ctx
	url_quota := 5
	rateLimitDuration := 24 * time.Hour

	rdb := database.CreateRedisClient()

	err := s.readJSON(w, r, &req)
	if err != nil {
		s.badRequestResponse(w, err)
		return
	}

	if req.CustomShortURL != "" {
		_, err := rdb.Get(ctx, req.CustomShortURL).Result()
		if err == redis.Nil {
			err = rdb.Set(context.Background(), req.CustomShortURL, req.URL, 0).Err()
			if err != nil {
				s.serverErrorResponse(w, err)
			}
			resp.CustomShortURL = req.CustomShortURL
		} else if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		} else {
			http.Error(w, "Short custom url already taken", http.StatusConflict)
			return
		}
	} else {
		resp.CustomShortURL = uuid.NewString()[:8]
		err = rdb.Set(context.Background(), resp.CustomShortURL, req.URL, 0).Err()
		if err != nil {
			s.serverErrorResponse(w, err)
		}
	}

	host := r.RemoteAddr
	rateLimit, err := rdb.Get(ctx, host).Int()
	if err == redis.Nil {
		rdb.Set(ctx, host, url_quota-1, rateLimitDuration).Err()
		rateLimit = url_quota - 1
	} else if err != nil {
		s.serverErrorResponse(w, err)
	} else {
		rateLimit64, err := rdb.Decr(ctx, host).Result()
		if err != nil {
			s.serverErrorResponse(w, err)
		}
		rateLimit = int(rateLimit64)
	}

	if rateLimit < 0 {
		ttl, err := rdb.TTL(ctx, host).Result()
		if err != nil {
			s.serverErrorResponse(w, err)
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

	resp.URL = req.URL
	resp.RateLimit = rateLimit

	err = s.writeJSON(w, http.StatusCreated, resp, nil)
	if err != nil {
		s.serverErrorResponse(w, err)
	}

}
