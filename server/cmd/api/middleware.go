package main

import (
	"log/slog"
	"net/http"
)

func (s *Server) Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		s.logger.Info("",
			slog.String("method", r.Method),
			slog.String("url", r.URL.String()),
			slog.String("remote_addr", r.RemoteAddr),
		)

		next.ServeHTTP(w, r)
	})

}
