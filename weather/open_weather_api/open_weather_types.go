package openweatherapi

// -- section: open weather API current
// https://openweathermap.org/current
type owCurrent struct {
	Weather []owCurrentWeather `json:"weather"`
	Main    owCurrentMain      `json:"main"`
}

type owCurrentWeather struct {
	Description string `json:"description"`
	Icon        string `json:"icon"`
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
	List    []owForecastData `json:"list"`
	Sunset  int64            `json:"sunrise"`
	Sunrise int64            `json:"sunset"`
}

type owForecastData struct {
	Timestamp int64                 `json:"dt"`
	Main      owForecastDataMain    `json:"main"`
	Weather   []owForecastDataWeather `json:"weather"`
	Pop       float64               `json:"pop"`
}

type owForecastDataMain struct {
	MinTemp float64 `json:"temp_min"`
	MaxTemp float64 `json:"temp_max"`
}

type owForecastDataWeather struct {
	Icon string `json:"icon"`
}
