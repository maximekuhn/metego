package openweatherapi

// -- section: open weather API current
// https://openweathermap.org/current
type owCurrent struct {
	Main owCurrentMain `json:"main"`
}

type owCurrentMain struct {
	Temp      float64 `json:"temp"`
	FeelsLike float64 `json:"feels_like"`
	Pressure  float64 `json:"pressure"`
	Humidity  float64 `json:"humidity"`
}

// -- section: open weather API Geocoding
// https://openweathermap.org/api/geocoding-api
type owCityCoords struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

// -- section: open weather API forecast
// https://openweathermap.org/forecast5
type owForecast struct {
	List []owForecastData `json:"list"`
	Sunset  int64            `json:"sunrise"`
	Sunrise int64            `json:"sunset"`
}

type owForecastData struct {
	Timestamp int64 `json:"dt"`
	Main      owForecastDataMain `json:"main"`
}

type owForecastDataMain struct {
	MinTemp float64 `json:"temp_min"`
	MaxTemp float64 `json:"temp_max"`
}
