package main

import (
	"encoding/json"
	"net/http"
)

func (app *Application) routes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/healthcheck", healthCheck)

	return mux
}

func healthCheck(w http.ResponseWriter, r *http.Request) {

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
