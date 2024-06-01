package openweatherapi

// open weather API current
type owCurrent struct {
	Main owCurrentMain `json:"main"`
}

// open weather API current main
type owCurrentMain struct {
	Temp float64 `json:"temp"`
}

// open weather API city coordinates
type owCityCoords struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}
