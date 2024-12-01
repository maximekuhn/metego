package openweatherapi

import (
	"testing"

	"github.com/maximekuhn/metego/internal/weather"
)

func TestConvertIconOk(t *testing.T) {
	tests := []struct {
		icon     string
		expected weather.CurrentWeatherIcon
	}{
		{"01d", weather.IconClearSky},
		{"01n", weather.IconNightClearSky},
		{"02d", weather.IconFewClouds},
		{"02n", weather.IconNightFewClouds},
		{"03d", weather.IconScatteredClouds},
		{"03n", weather.IconNightScatteredClouds},
		{"04d", weather.IconBrokenClouds},
		{"04n", weather.IconNightBrokenClouds},
		{"09d", weather.IconShowerRain},
		{"09n", weather.IconNightShowerRain},
		{"10d", weather.IconRain},
		{"10n", weather.IconNightRain},
		{"11d", weather.IconThunderstorm},
		{"11n", weather.IconNightThunderstorm},
		{"13d", weather.IconSnow},
		{"13n", weather.IconNightSnow},
		{"50d", weather.IconMist},
		{"50n", weather.IconNightMist},
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
