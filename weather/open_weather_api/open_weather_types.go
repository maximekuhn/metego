package openweatherapi

// open weather API current
type owCurrent struct {
	Main owCurrentMain `json:"main"`
}

// open weather API current main
type owCurrentMain struct {
	Temp      float64 `json:"temp"`
	FeelsLike float64 `json:"feels_like"`
	Pressure  float64 `json:"pressure"`
	Humidity  float64 `json:"humidity"`
}

// open weather API city coordinates
type owCityCoords struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}
