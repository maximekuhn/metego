package server

import (
	"log/slog"
	"net/http"

	"github.com/maximekuhn/metego/server/views"
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

func (s *Server) weatherHandler(w http.ResponseWriter, r *http.Request) {
	city := r.PathValue("city")
	slog.Info("GET /weather/{city}", slog.String("city", city))
	weatherPage := views.Weather(city)

	w.Header().Add("Content-Type", "text/html")
	err := weatherPage.Render(r.Context(), w)
	if err != nil {
		slog.Error("failed to render Weather", slog.String("err_msg", err.Error()))
	}
}

func (s *Server) currentWeatherHandler(w http.ResponseWriter, r *http.Request) {
	city := r.URL.Query().Get("city")
	slog.Info("GET /api/weather/current", slog.String("city", city))
	currentWeather, err := s.state.fetcher.FetchCurrent(city)
	if err != nil {
		slog.Error("failed to get current weather", slog.String("err_msg", err.Error()))
		// TODO: check err, maybe it's the user's fault
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	currentWeatherComp := views.CurrentWeatherComp(city, currentWeather)
	w.Header().Add("Content-Type", "text/html")
	err = currentWeatherComp.Render(r.Context(), w)
	if err != nil {
		slog.Error("failed to render CurrentWeatherComp", slog.String("err_msg", err.Error()))
	}
}
