package server

import (
	"net/http"

	"github.com/maximekuhn/metego/weather"
)

type Server struct {
	state *state
}

func NewServer(fetcher weather.Fetcher) *Server {
	return &Server{
		state: NewState(fetcher),
	}
}

func (s *Server) Start() error {
	http.HandleFunc("GET /weather/{city}", s.weatherHandler)
	http.HandleFunc("GET /api/weather/current/", s.currentWeatherHandler)
	err := http.ListenAndServe(":9004", nil)
	return err
}
