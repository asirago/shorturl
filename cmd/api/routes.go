package main

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *application) routes() http.Handler {
	r := chi.NewRouter()

	r.Get("/healthcheck", app.healthCheck)

	r.Post("/shorten-url", app.shortenUrl)
	return r
}

func (app *application) healthCheck(w http.ResponseWriter, r *http.Request) {

	data := map[string]string{
		"version":     "1.0.0",
		"status":      "available",
		"environment": "development",
	}

	json, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(json)
}
