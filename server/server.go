package server

import (
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
	http.ListenAndServe(":9004", nil)
	return nil
}

func (s *Server) weatherHandler(w http.ResponseWriter, r *http.Request) {
    city := r.PathValue("city")
    currentWeather, err := s.state.fetcher.FetchCurrent(city)
    if err != nil {
        http.Error(w, "Internal Server Error", 500)
        return
    }

    weatherPage := views.CurrentWeather(city, currentWeather)

    w.Header().Add("Content-Type", "text/html")
    weatherPage.Render(r.Context(), w)
    


}
