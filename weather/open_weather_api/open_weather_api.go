package openweatherapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/maximekuhn/metego/weather"
)

type OpenWeatherFetcher struct {
	apiKey string
	cities map[string]*owCityCoords
}

func NewOpenWeatherFetcher(apiKey string) *OpenWeatherFetcher {
	return &OpenWeatherFetcher{
		apiKey: apiKey,
		cities: map[string]*owCityCoords{},
	}
}

func (f *OpenWeatherFetcher) FetchCurrent(city string) (*weather.CurrentWeather, error) {
	curr, err := f.fetchCurrent(city)
	if err != nil {
		return nil, err
	}

	current := weather.CurrentWeather{
		Temp: curr.Main.Temp,
	}

	return &current, nil
}

// fetch the current weather for a given city
func (f *OpenWeatherFetcher) fetchCurrent(city string) (*owCurrent, error) {
	// open weather requires city coordinates and not city name
	coords, ok := f.cities[city]
	if !ok {
		cityCoords, err := f.fetchCityCoords(city)
		if err != nil {
			return nil, err
		}
		f.cities[city] = cityCoords
		coords = cityCoords
	}

	url := fmt.Sprintf(
		"https://api.openweathermap.org/data/2.5/weather?lat=%f&lon=%f&appid=%s&units=metric&lang=fr",
		coords.Lat,
		coords.Lon,
		f.apiKey,
	)

	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != 200 {
		return nil, errors.New("failed to GET current weather")
	}

	decoder := json.NewDecoder(res.Body)
	var current owCurrent
	err = decoder.Decode(&current)
	if err != nil {
		return nil, err
	}

	return &current, nil
}

// fetch city coordinates
// if coords are not found, then an error is returned
func (f *OpenWeatherFetcher) fetchCityCoords(city string) (*owCityCoords, error) {
	url := fmt.Sprintf(
		"https://api.openweathermap.org/geo/1.0/direct?q=%s&appid=%s",
		city,
		f.apiKey,
	)

	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != 200 {
		return nil, errors.New("failed to GET city coords")
	}

	decoder := json.NewDecoder(res.Body)
	var cityCoords []owCityCoords
	err = decoder.Decode(&cityCoords)
	if err != nil {
		return nil, err
	}

	// get only first coords
	if len(cityCoords) == 0 {
		return nil, errors.New("no coordinates found")
	}

	coords := cityCoords[0]

	return &coords, nil
}
