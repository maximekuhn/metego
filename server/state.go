package server

import (
	"github.com/maximekuhn/metego/calendar"
	"github.com/maximekuhn/metego/weather"
)

type state struct {
	fetcher weather.Fetcher
	storage calendar.BirthdayStorage
}

func NewState(fetcher weather.Fetcher, storage calendar.BirthdayStorage) *state {
	return &state{
		fetcher: fetcher,
		storage: storage,
	}
}
