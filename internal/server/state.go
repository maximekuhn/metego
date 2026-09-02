package server

import (
	"sync"

	"github.com/maximekuhn/metego/internal/calendar"
	"github.com/maximekuhn/metego/internal/weather"
)

type state struct {
	weatherProviders []WeatherProvider
	currProviderMu   sync.Mutex
	currProviderIdx  int

	bdaysStorage    calendar.BirhtdayStorage
	aptsStorage     calendar.AppointmentStorage
	namedaysStorage calendar.NamedayStorage
	cities          []string
}

func NewState(
	weatherProviders []WeatherProvider,
	bdaysStorage calendar.BirhtdayStorage,
	aptsStorage calendar.AppointmentStorage,
	namedayStorage calendar.NamedayStorage,
	cities []string,
) *state {
	if len(weatherProviders) == 0 {
		panic("server: state must have at least 1 weather provider, 0 provided")
	}

	return &state{
		weatherProviders: weatherProviders,
		currProviderMu:   sync.Mutex{},
		currProviderIdx:  0,

		bdaysStorage:    bdaysStorage,
		aptsStorage:     aptsStorage,
		namedaysStorage: namedayStorage,
		cities:          cities,
	}
}

// nextCity returns the next city to display to the user
// if there is no next city (cities list is empty), an empty string is returned
func (s *state) nextCity(currentCity string) string {
	if len(s.cities) == 0 {
		return ""
	}

	cityIdx := -1
	for idx, city := range s.cities {
		if city == currentCity {
			cityIdx = idx
			break
		}
	}

	// if the current city isn't in the cities list,
	// we default to the first one
	if cityIdx == -1 {
		return s.cities[0]
	}

	if cityIdx < len(s.cities)-1 {
		return s.cities[cityIdx+1]
	}

	return s.cities[0]

}

func (s *state) currentWeatherProvider() WeatherProvider {
	s.currProviderMu.Lock()
	defer s.currProviderMu.Unlock()
	return s.weatherProviders[s.currProviderIdx]
}

func (s *state) updateWeatherProvider() {
	s.currProviderMu.Lock()
	defer s.currProviderMu.Unlock()
	s.currProviderIdx = (s.currProviderIdx + 1) % len(s.weatherProviders)
}

type WeatherProvider struct {
	weather.Fetcher
	SourceName string
	URL        string
}
