package openweatherapi

import (
	"testing"

	"github.com/maximekuhn/metego/weather"
)

func TestConvertIconOk(t *testing.T) {
	tests := []struct {
		icon     string
		expected weather.CurrentWeatherIcon
	}{
		{"01d", weather.IconClearSky},
		{"01n", weather.IconClearSky},
		{"02d", weather.IconFewClouds},
		{"02n", weather.IconFewClouds},
		{"03d", weather.IconScatteredClouds},
		{"03n", weather.IconScatteredClouds},
		{"04d", weather.IconBrokenClouds},
		{"04n", weather.IconBrokenClouds},
		{"09d", weather.IconShowerRain},
		{"09n", weather.IconShowerRain},
		{"10d", weather.IconRain},
		{"10n", weather.IconRain},
		{"11d", weather.IconThunderstorm},
		{"11n", weather.IconThunderstorm},
		{"13d", weather.IconSnow},
		{"13n", weather.IconSnow},
		{"50d", weather.IconMist},
		{"50n", weather.IconMist},
	}

	for _, test := range tests {
		t.Run(test.icon, func(t *testing.T) {
			actual, err := toWeatherIcon(test.icon)
			if err != nil {
				t.Errorf("got an unexpected error: %s", err)
			}

			if actual != test.expected {
				t.Errorf("want %d got %d", test.expected, actual)
			}
		})
	}
}
