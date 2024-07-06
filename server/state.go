package server

import (
	"github.com/maximekuhn/metego/calendar"
	"github.com/maximekuhn/metego/weather"
)

type state struct {
	fetcher      weather.Fetcher
	bdaysStorage calendar.BirhtdayStorage
	aptsStorage  calendar.AppointmentStorage
}

func NewState(fetcher weather.Fetcher, bdaysStorage calendar.BirhtdayStorage, aptsStorage calendar.AppointmentStorage) *state {
	return &state{
		fetcher:      fetcher,
		bdaysStorage: bdaysStorage,
		aptsStorage:  aptsStorage,
	}
}
