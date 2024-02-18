package main

import "net/http"

func (s *Server) errorJSON(w http.ResponseWriter, status int, message string) {

	data := map[string]string{
		"error": message,
	}

	err := s.writeJSON(w, status, data, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	}
}

func (s *Server) serverErrorResponse(w http.ResponseWriter, err error) {
	s.logger.Error("", "error", err.Error())

	s.errorJSON(
		w,
		http.StatusInternalServerError,
		"server encountered a problem and could not process your request",
	)
}

func (s *Server) badRequestResponse(w http.ResponseWriter, err error) {
	s.errorJSON(
		w,
		http.StatusBadRequest,
		err.Error(),
	)
}

func (s *Server) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	s.errorJSON(w, http.StatusNotFound, "requested resource could not be found")
}
