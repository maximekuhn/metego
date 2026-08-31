package openmeteoapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sync"

	"github.com/maximekuhn/metego/internal/weather"
)

type OpenMeteoFetcher struct {
	coordsMu sync.Mutex
	coords   map[string] /* city name */ coords
}

func NewOpenMeteoFetcher() *OpenMeteoFetcher {
	return &OpenMeteoFetcher{
		coordsMu: sync.Mutex{},
		coords:   map[string]coords{},
	}
}

func (f *OpenMeteoFetcher) FetchCurrent(city string) (*weather.CurrentWeather, error) {
	coords, err := f.getCityCoords(city)
	if err != nil {
		return nil, fmt.Errorf("openmeteoapi: could not fetch coords: %w", err)
	}

	url := fmt.Sprintf(
		"https://api.open-meteo.com/v1/forecast?latitude=%.2f&longitude=%.2f&daily=sunrise,sunset,temperature_2m_max,temperature_2m_min,precipitation_probability_max&hourly=wind_speed_10m,surface_pressure,weather_code,temperature_2m&current=temperature_2m,wind_speed_10m,weather_code,surface_pressure,relative_humidity_2m,is_day&timezone=%s&timeformat=unixtime",
		coords.lat,
		coords.lon,
		url.QueryEscape("Europe/Berlin"),
	)

	res, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("openmeteoapi: could not fetch weather: %w", err)
	}
	defer res.Body.Close()

	var resp weatherResponse
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		return nil, fmt.Errorf("openmeteoapi: failed to deserialize response: %w", err)
	}

	if len(resp.Daily.Sunset) < 1 || len(resp.Daily.Sunrise) < 1 {
		return nil, fmt.Errorf("openmeteoapi: no Daily sunset or sunrise")
	}

	return &weather.CurrentWeather{
		Temp:        resp.Current.Temperature2M,
		Pressure:    resp.Current.SurfacePressure,
		Humidity:    resp.Current.RelativeHumidity2M,
		WindSpeed:   resp.Current.WindSpeed10M,
		Description: toDescription(resp.Current.WeatherCode),
		Icon:        toIcon(resp.Current.WeatherCode, resp.Current.IsDay == 1),
		Sunset:      resp.Daily.Sunset[0],
		Sunrise:     resp.Daily.Sunrise[0],
	}, nil
}

func (f *OpenMeteoFetcher) FetchForecast(city string, days int) ([]*weather.ForecastWeather, error) {
	panic("TODO")
}

func (f *OpenMeteoFetcher) fetchCityCoords(cityName string) (*coords, error) {
	url := fmt.Sprintf(
		"https://geocoding-api.open-meteo.com/v1/search?name=%s&count=1&language=en&format=json",
		url.QueryEscape(cityName),
	)

	res, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get OpenMeteo/GeoCoding results for city: '%s': %w", cityName, err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("OpenMeteo/Geocoding: expected 200 OK, got %s", res.Status)
	}

	type response struct {
		Results []struct {
			Lat float64 `json:"latitude"`
			Lon float64 `json:"longitude"`
		} `json:"results"`
	}

	var resp response
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		return nil, fmt.Errorf("OpenMeteo/Geocoding: failed to deserialize response: %w", err)
	}

	if len(resp.Results) != 1 {
		return nil, fmt.Errorf("OpenMeteo/Geocoding: no result for city '%s'", cityName)
	}

	return &coords{
		lat: resp.Results[0].Lat,
		lon: resp.Results[0].Lon,
	}, nil
}

func (f *OpenMeteoFetcher) getCityCoords(cityName string) (*coords, error) {
	f.coordsMu.Lock()
	defer f.coordsMu.Unlock()

	if c, found := f.coords[cityName]; found {
		return &c, nil
	}

	c, err := f.fetchCityCoords(cityName)
	if err != nil {
		return nil, err
	}
	f.coords[cityName] = *c
	return c, nil
}

func toIcon(weatherCode int, isDay bool) weather.CurrentWeatherIcon {
	var icon weather.CurrentWeatherIcon

	switch weatherCode {
	case 0:
		icon = weather.IconClearSky

	case 1:
		icon = weather.IconFewClouds
	case 2:
		icon = weather.IconScatteredClouds
	case 3:
		icon = weather.IconBrokenClouds

	case 45, 48:
		icon = weather.IconMist

	case 51, 53, 55:
		icon = weather.IconRain
	case 56, 57:
		icon = weather.IconRain

	case 61, 63, 65:
		icon = weather.IconRain
	case 66, 67:
		icon = weather.IconRain

	case 71, 73, 75, 77:
		icon = weather.IconSnow

	case 80, 81, 82:
		icon = weather.IconShowerRain

	case 85, 86:
		icon = weather.IconSnow

	case 95, 96, 99:
		icon = weather.IconThunderstorm

	default:
		icon = weather.IconClearSky
	}

	if isDay {
		return icon
	}

	switch icon {
	case weather.IconClearSky:
		return weather.IconNightClearSky
	case weather.IconFewClouds:
		return weather.IconNightFewClouds
	case weather.IconScatteredClouds:
		return weather.IconNightScatteredClouds
	case weather.IconBrokenClouds:
		return weather.IconNightBrokenClouds
	case weather.IconShowerRain:
		return weather.IconNightShowerRain
	case weather.IconRain:
		return weather.IconNightRain
	case weather.IconThunderstorm:
		return weather.IconNightThunderstorm
	case weather.IconSnow:
		return weather.IconNightSnow
	case weather.IconMist:
		return weather.IconNightMist
	default:
		return weather.IconNightClearSky
	}
}

func toDescription(weatherCode int) string {
	switch weatherCode {
	case 0:
		return "Ciel dégagé"

	case 1:
		return "Principalement dégagé"
	case 2:
		return "Partiellement nuageux"
	case 3:
		return "Couvert"

	case 45:
		return "Brouillard"
	case 48:
		return "Brouillard givrant"

	case 51:
		return "Bruine légère"
	case 53:
		return "Bruine modérée"
	case 55:
		return "Bruine forte"

	case 56:
		return "Bruine verglaçante légère"
	case 57:
		return "Bruine verglaçante forte"

	case 61:
		return "Pluie faible"
	case 63:
		return "Pluie modérée"
	case 65:
		return "Pluie forte"

	case 66:
		return "Pluie verglaçante légère"
	case 67:
		return "Pluie verglaçante forte"

	case 71:
		return "Chutes de neige faibles"
	case 73:
		return "Chutes de neige modérées"
	case 75:
		return "Chutes de neige fortes"

	case 77:
		return "Grains de neige"

	case 80:
		return "Averses de pluie faibles"
	case 81:
		return "Averses de pluie modérées"
	case 82:
		return "Averses de pluie violentes"

	case 85:
		return "Averses de neige faibles"
	case 86:
		return "Averses de neige fortes"

	case 95:
		return "Orage faible ou modéré"

	case 96:
		return "Orage avec grêle faible"
	case 99:
		return "Orage avec grêle forte"

	default:
		return "Conditions météorologiques inconnues"
	}
}

type coords struct {
	lat float64
	lon float64
}

type weatherResponse struct {
	Latitude             float64 `json:"latitude"`
	Longitude            float64 `json:"longitude"`
	GenerationTimeMS     float64 `json:"generationtime_ms"`
	UTCOffsetSeconds     int     `json:"utc_offset_seconds"`
	Timezone             string  `json:"timezone"`
	TimezoneAbbreviation string  `json:"timezone_abbreviation"`
	Elevation            float64 `json:"elevation"`

	CurrentUnits currentUnits `json:"current_units"`
	Current      current      `json:"current"`

	HourlyUnits hourlyUnits `json:"hourly_units"`
	Hourly      hourly      `json:"hourly"`

	DailyUnits dailyUnits `json:"daily_units"`
	Daily      daily      `json:"daily"`
}

type currentUnits struct {
	Time            string `json:"time"`
	Interval        string `json:"interval"`
	Temperature2M   string `json:"temperature_2m"`
	WindSpeed10M    string `json:"wind_speed_10m"`
	WeatherCode     string `json:"weather_code"`
	SurfacePressure string `json:"surface_pressure"`
}

type current struct {
	Interval           int     `json:"interval"`
	Temperature2M      float64 `json:"temperature_2m"`
	WindSpeed10M       float64 `json:"wind_speed_10m"`
	WeatherCode        int     `json:"weather_code"`
	SurfacePressure    float64 `json:"surface_pressure"`
	RelativeHumidity2M float64 `json:"relative_humidity_2m"`
	IsDay              int     `json:"is_day"`
}

type hourlyUnits struct {
	Time            string `json:"time"`
	WindSpeed10M    string `json:"wind_speed_10m"`
	SurfacePressure string `json:"surface_pressure"`
	WeatherCode     string `json:"weather_code"`
	Temperature2M   string `json:"temperature_2m"`
}

type hourly struct {
	WindSpeed10M    []float64 `json:"wind_speed_10m"`
	SurfacePressure []float64 `json:"surface_pressure"`
	WeatherCode     []int     `json:"weather_code"`
	Temperature2M   []float64 `json:"temperature_2m"`
}

type dailyUnits struct {
	Time                     string `json:"time"`
	Sunrise                  string `json:"sunrise"`
	Sunset                   string `json:"sunset"`
	Temperature2MMax         string `json:"temperature_2m_max"`
	Temperature2MMin         string `json:"temperature_2m_min"`
	PrecipitationProbability string `json:"precipitation_probability_max"`
}

type daily struct {
	Sunrise                  []int64   `json:"sunrise"`
	Sunset                   []int64   `json:"sunset"`
	Temperature2MMax         []float64 `json:"temperature_2m_max"`
	Temperature2MMin         []float64 `json:"temperature_2m_min"`
	PrecipitationProbability []int     `json:"precipitation_probability_max"`
}
