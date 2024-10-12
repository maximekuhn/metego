package server

import (
	"github.com/maximekuhn/metego/calendar"
	"github.com/maximekuhn/metego/weather"
)

type state struct {
	fetcher      weather.Fetcher
	bdaysStorage calendar.BirhtdayStorage
	aptsStorage  calendar.AppointmentStorage
	cities       []string
}

func NewState(
	fetcher weather.Fetcher,
	bdaysStorage calendar.BirhtdayStorage,
	aptsStorage calendar.AppointmentStorage,
	cities []string,
) *state {
	return &state{
		fetcher:      fetcher,
		bdaysStorage: bdaysStorage,
		aptsStorage:  aptsStorage,
		cities:       cities,
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
