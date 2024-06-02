package weather

type Fetcher interface {
	FetchCurrent(city string) (*CurrentWeather, error)

	// fetch `days` of forecast data
	// if it's not possible to get all days, no error is returned but the array doesn't contain all days
	FetchForecast(city string, days int) ([]*ForecastWeather, error)
}
