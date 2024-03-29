package main

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (s *Server) routes() http.Handler {
	r := chi.NewRouter()

	r.NotFound(s.notFoundResponse)

	r.Use(s.Logger)
	r.Get("/healthcheck", s.healthCheck)

	r.Post("/shorten-url", s.shortenUrl)
	r.Get("/{url}", s.resolveUrl)
	return r
}

func (s *Server) healthCheck(w http.ResponseWriter, r *http.Request) {

	data := map[string]string{
		"version":     "1.0.0",
		"status":      "available",
		"environment": s.cfg.Environment,
	}

	json, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(json)
}
