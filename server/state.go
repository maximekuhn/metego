package server

import (
	"github.com/maximekuhn/metego/calendar"
	"github.com/maximekuhn/metego/weather"
)

type state struct {
	fetcher      weather.Fetcher
	bdaysStorage calendar.BirhtdayStorage
}

func NewState(fetcher weather.Fetcher, bdaysStorage calendar.BirhtdayStorage) *state {
	return &state{
		fetcher:      fetcher,
		bdaysStorage: bdaysStorage,
	}
}
