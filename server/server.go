package server

import (
	"net/http"

	"github.com/maximekuhn/metego/calendar"
	"github.com/maximekuhn/metego/weather"
)

type Server struct {
	state *state
}

func NewServer(fetcher weather.Fetcher, storage calendar.BirthdayStorage) *Server {
	return &Server{
		state: NewState(fetcher, storage),
	}
}

func (s *Server) Start() error {
	// index routes
	http.HandleFunc("GET /", s.handleRoot)

	// weather routes
	http.HandleFunc("GET /weather/{city}", s.weatherHandler)
	http.HandleFunc("GET /api/weather/current/", s.currentWeatherHandler)
	http.HandleFunc("GET /api/weather/forecast/", s.handleGetForecastWeather)

	// birthdays routes
	http.HandleFunc("GET /birthdays", s.birthdaysHandler)
	http.HandleFunc("GET /api/birthdays", s.handleGetTodayBirthdays)
	http.HandleFunc("POST /api/birthdays", s.handleCreateBirthday)

	// static files
	http.Handle("GET /static/", http.StripPrefix("/static", http.FileServer(http.Dir("./static"))))

	// TODO: get this from conf
	err := http.ListenAndServe(":9004", nil)
	return err
}
