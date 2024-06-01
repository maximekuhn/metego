//go:build integration

package openweatherapi

import (
	"errors"
	"math"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func TestFetchCityCoords(t *testing.T) {
	tests := []struct {
		city        string
		expectedLat float64
		expectedLon float64
	}{
		{
			"Paris",
			48.8,
			2.3,
		},
		{
			"Salt Like City",
			40.7,
			-111.8,
		},
	}

	for _, test := range tests {
		t.Run(test.city, func(t *testing.T) {
			fetcher, err := setup()
			if err != nil {
				t.Fatalf("failed to setup fetcher: %s", err)
			}

			city := "Paris"
			expectedLat := 48.8
			expectedLon := 2.3

			coords, err := fetcher.fetchCityCoords(city)
			if err != nil {
				t.Fatalf("failed to get city coords: %s", err)
			}
			if nearlyEquals(expectedLat, coords.Lat, 0.01) {
				t.Errorf("lat: want %.1f got %.1f", expectedLat, coords.Lat)
			}
			if nearlyEquals(expectedLon, coords.Lon, 0.01) {
				t.Errorf("lon: want %.1f got %.1f", expectedLon, coords.Lon)
			}
		})
	}
}

func TestFetchCityCoordsImaginary(t *testing.T) {
	fetcher, err := setup()
	if err != nil {
		t.Fatalf("failed to setup fetcher: %s", err)
	}

	city := "Bikini Bottom"
	_, err = fetcher.fetchCityCoords(city)
	if err == nil {
		t.Fatalf("got city coords but expected error")
	}

	expectedErrMsg := "no coordinates found"
	actualErrMsg := err.Error()
	if actualErrMsg != expectedErrMsg {
		t.Fatalf("want %s got %s", expectedErrMsg, actualErrMsg)
	}
}

func TestFetchCurrentWeather(t *testing.T) {
	fetcher, err := setup()
	if err != nil {
		t.Fatalf("failed to setup fetcher: %s", err)
	}

	city := "Paris"
	current, err := fetcher.fetchCurrent(city)
	if err != nil {
		t.Fatalf("expected current weather got error: %s", err)
	}

	if current == nil {
		t.Fatalf("current weather is nil")
	}
}

func TestFetchCurrentWeatherUpdatesCache(t *testing.T) {
	tests := []struct {
		city string
	}{
		{
			"Paris",
		},
		{
			"Berlin",
		},
		{
			"Stockholm",
		},
	}

	for _, test := range tests {
		t.Run(test.city, func(t *testing.T) {
			fetcher, err := setup()
			if err != nil {
				t.Fatalf("failed to setup fetcher: %s", err)
			}

			// check cache miss for city name
			_, ok := fetcher.cities[test.city]
			if ok {
				t.Fatalf("go cache hit expected cache miss")
			}

			_, err = fetcher.fetchCurrent(test.city)
			if err != nil {
				t.Fatalf("failed to fetch current: %s", err)
			}

			// check cache hit
			_, ok = fetcher.cities[test.city]
			if !ok {
				t.Fatalf("go cache miss expected cache hit")
			}
		})
	}
}

func nearlyEquals(a, b, precision float64) bool {
	return math.Abs(a-b) < precision
}

func setup() (*openWeatherFetcher, error) {
	err := godotenv.Load("../../.env.integration")
	if err != nil {
		return nil, err
	}

	openWeatherApiKey := os.Getenv("OPEN_WEATHER_API_KEY")
	if openWeatherApiKey == "" {
		return nil, errors.New("OPEN_WEATHER_API_KEY is not set")
	}

	return newOpenWeatherFetcher(openWeatherApiKey), nil
}
