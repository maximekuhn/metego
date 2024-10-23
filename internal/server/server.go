package server

import (
	"net/http"

	"github.com/maximekuhn/metego/internal/calendar"
	"github.com/maximekuhn/metego/internal/weather"
)

type Server struct {
	state *state
}

// NewServer creates a new server :)
//
// The parameter `cities` is a list of pre-defined cities that will be used
// when the user clicks on the city name in the menu bar.
//
// For instance, cities can be: ["Paris", "Berlin"].
// If the app is running and the user click on "Paris", it will switch to "Berlin"
func NewServer(
	fetcher weather.Fetcher,
	bdaysStorage calendar.BirhtdayStorage,
	aptsStorage calendar.AppointmentStorage,
	cities []string,
) *Server {
	return &Server{
		state: NewState(fetcher, bdaysStorage, aptsStorage, cities),
	}
}

func (s *Server) Start() error {
	// index routes
	http.HandleFunc("GET /", s.handleRoot)
	http.HandleFunc("GET /admin", s.handleAdmin)

	// weather routes
	http.HandleFunc("GET /weather/{city}", s.weatherHandler)
	http.HandleFunc("GET /api/weather/current/", s.currentWeatherHandler)
	http.HandleFunc("GET /api/weather/current/metrics", s.currentMetricsWeatherHandler)
	http.HandleFunc("GET /api/weather/forecast/", s.handleGetForecastWeather)

	// birthdays routes
	http.HandleFunc("GET /birthdays", s.birthdaysHandler)
	http.HandleFunc("GET /api/birthdays", s.handleGetTodayBirthdays)
	http.HandleFunc("POST /api/birthdays", s.handleCreateBirthday)

	// appointments routes
	http.HandleFunc("GET /appointments", s.appointmentsHandler)
	http.HandleFunc("POST /api/appointments", s.handleCreateAppointment)

	// cities route
	http.HandleFunc("GET /city/next", s.changeCityHandler)

	// static files
	http.Handle("GET /static/", http.StripPrefix("/static", http.FileServer(http.Dir("./static"))))

	// TODO: get this from conf
	err := http.ListenAndServe(":9004", nil)
	return err
}
