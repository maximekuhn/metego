package weather

type Fetcher interface {
	FetchCurrent(city string) (*CurrentWeather, error)
}
