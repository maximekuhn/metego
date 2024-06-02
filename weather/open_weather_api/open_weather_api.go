package openweatherapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/maximekuhn/metego/weather"
)

type OpenWeatherFetcher struct {
	apiKey string

	// TODO: prevent data races (mutex or rw lock)
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
		Temp:      curr.Main.Temp,
		FeelsLike: curr.Main.FeelsLike,
		Pressure:  curr.Main.Pressure,
		Humidity:  curr.Main.Humidity,
	}

	return &current, nil
}

func (f *OpenWeatherFetcher) FetchForecast(city string, days int) ([]*weather.ForecastWeather, error) {
	fcs, err := f.fetchForecast(city)
	if err != nil {
		return nil, err
	}

	// get noon forecast for each day (UTC)
	// ignore today
	today := time.Now().Day()
	forecasts := make([]*weather.ForecastWeather, 0)
	for _, f := range fcs.List {
		time := time.Unix(f.Timestamp, 0)

		if time.Day() == today {
			continue
		}

		if time.UTC().Hour() != 12 {
			continue
		}

		// TODO: fix highest and lowest (they are the same)

		forecast := &weather.ForecastWeather{
			Date:        time,
			HighestTemp: f.Main.MaxTemp,
			LowestTemp:  f.Main.MinTemp,
			Pop:         0.0, // TODO
		}

		forecasts = append(forecasts, forecast)
	}

	//return forecasts, nil
    /// XXX: is this dangerous ?
	return forecasts[0:days], nil
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

// fetch forecast weather for the next 5 days for a given city
func (f *OpenWeatherFetcher) fetchForecast(city string) (*owForecast, error) {
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
		"https://api.openweathermap.org/data/2.5/forecast?lat=%f&lon=%f&appid=%s&units=metric&lang=fr",
		coords.Lat,
		coords.Lon,
		f.apiKey,
	)
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != 200 {
		return nil, errors.New("failed to GET forecast weather")
	}

	decoder := json.NewDecoder(res.Body)
	var forecasts *owForecast
	err = decoder.Decode(&forecasts)
	if err != nil {
		return nil, err
	}

	return forecasts, nil
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
