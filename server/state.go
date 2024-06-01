package server

import "github.com/maximekuhn/metego/weather"

type state struct {
    fetcher weather.Fetcher
}

func NewState(fetcher weather.Fetcher) *state {
    return &state{
    	fetcher: fetcher,
    }
}
