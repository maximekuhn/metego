package server

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/maximekuhn/metego/internal/server/views"
)

// GET /weather/{city}
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

// GET /api/weather/current?city=...
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

	currentWeatherComp := views.CurrentWeatherComp(currentWeather)
	w.Header().Add("Content-Type", "text/html")
	err = currentWeatherComp.Render(r.Context(), w)
	if err != nil {
		slog.Error("failed to render CurrentWeatherComp", slog.String("err_msg", err.Error()))
	}
}

// GET /api/weather/current/metrics?city=...
func (s *Server) currentMetricsWeatherHandler(w http.ResponseWriter, r *http.Request) {
	city := r.URL.Query().Get("city")
	slog.Info("GET /api/weather/current/metrics", slog.String("city", city))
	currentWeather, err := s.state.fetcher.FetchCurrent(city)
	if err != nil {
		slog.Error("failed to get current weather", slog.String("err_msg", err.Error()))
		// TODO: check err, maybe it's the user's fault
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	todayMetricsComp := views.TodayMetricsComp(currentWeather)
	w.Header().Add("Content-Type", "text/html")
	err = todayMetricsComp.Render(r.Context(), w)
	if err != nil {
		slog.Error("failed to render TodayMetricsComp", slog.String("err_msg", err.Error()))
	}
}

// GET /api/weather/forecast?city=...&days=...
func (s *Server) handleGetForecastWeather(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	city := queryParams.Get("city")
	days := queryParams.Get("days")
	slog.Info(
		"GET /api/weather/forecast",
		slog.String("city", city),
		slog.String("days", days),
	)

	d, err := strconv.ParseInt(days, 10, 64)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	forecasts, err := s.state.fetcher.FetchForecast(city, int(d))
	if err != nil {
		slog.Error("failed to get forecast weather", slog.String("err_msg", err.Error()))
		// TODO: check err, maybe it's the user's fault
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	slog.Info("got forecasts", slog.Int("days", len(forecasts)))

	forecastWeatherComp := views.ForecastWeatherComp(forecasts)
	w.Header().Add("Content-Type", "text/html")
	err = forecastWeatherComp.Render(r.Context(), w)
	if err != nil {
		slog.Error("failed to render ForecastWeatherComp", slog.String("err_msg", err.Error()))
	}
}
