package server

import (
	"fmt"
	"net/http"
)

// GET /city/next?current={city}
func (s *Server) changeCityHandler(w http.ResponseWriter, r *http.Request) {
	currentCity := r.URL.Query().Get("current")
	fmt.Printf("GET /city/next?current=%s\n", currentCity)
	nextCity := s.state.nextCity(currentCity)
	if nextCity == "" {
		nextCity = currentCity
	}
	redirectURL := fmt.Sprintf("/weather?city=%s", nextCity)

	// TODO: maybe return a redirect but then it need to be handled in the front
	// and I don't know/want to know how for now.
	w.Header().Add("Content-Type", "text/plain")
	_, err := w.Write([]byte(redirectURL))
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
	}
}
