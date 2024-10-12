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
