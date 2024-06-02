package weather

import "time"

type ForecastWeather struct {
	// UTC
	Date time.Time

	HighestTemp float64
	LowestTemp  float64

	// percentage of precipitation
	Pop float64
}
