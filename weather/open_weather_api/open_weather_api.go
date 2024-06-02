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

	forecastsByDay := make(map[int]*weather.ForecastWeather, 0)
	today := time.Now().Day()
	for _, f := range fcs.List {
		time := time.Unix(f.Timestamp, 0)
		day := time.Day()

		// ignore today forecast
		if day == today {
			continue
		}

		forecastDay, ok := forecastsByDay[day]
		if !ok {
			forecastsByDay[day] = &weather.ForecastWeather{
				Date:        time,
				HighestTemp: f.Main.MaxTemp,
				LowestTemp:  f.Main.MinTemp,
				Pop:         0, // TODO
			}

			forecastDay = forecastsByDay[day]
		}

		if f.Main.MinTemp < forecastDay.LowestTemp {
			forecastDay.LowestTemp = f.Main.MinTemp
		}
		if f.Main.MaxTemp > forecastDay.HighestTemp {
			forecastDay.HighestTemp = f.Main.MaxTemp
		}
	}

	forecasts := make([]*weather.ForecastWeather, 0)
	for _, f := range forecastsByDay {
		forecasts = append(forecasts, f)
	}

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
