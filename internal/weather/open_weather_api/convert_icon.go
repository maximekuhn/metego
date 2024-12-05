package openweatherapi

import (
	"errors"
	"fmt"

	"github.com/maximekuhn/metego/internal/weather"
)

// try to convert the icon string from open weather API to a custom type
//
// if an error si returned ,the icon is not valid
func toWeatherIcon(icon string) (weather.CurrentWeatherIcon, error) {
	if len(icon) != 3 {
		return 0, errors.New("unknown icon type")
	}

	last := icon[2]
	if last == 'd' {
		return dayIcon(icon)
	} else if last == 'n' {
		return nightIcon(icon)
	}

	return weather.IconClearSky, fmt.Errorf("unknown icon: %s", icon)
}

func dayIcon(icon string) (weather.CurrentWeatherIcon, error) {
	switch icon[:2] {
	case "01":
		return weather.IconClearSky, nil
	case "02":
		return weather.IconFewClouds, nil
	case "03":
		return weather.IconScatteredClouds, nil
	case "04":
		return weather.IconBrokenClouds, nil
	case "09":
		return weather.IconShowerRain, nil
	case "10":
		return weather.IconRain, nil
	case "11":
		return weather.IconThunderstorm, nil
	case "13":
		return weather.IconSnow, nil
	case "50":
		return weather.IconMist, nil
	}
	return weather.IconClearSky, fmt.Errorf("dayIcon(): unknown icon %s", icon)
}

func nightIcon(icon string) (weather.CurrentWeatherIcon, error) {
	switch icon[:2] {
	case "01":
		return weather.IconNightClearSky, nil
	case "02":
		return weather.IconNightFewClouds, nil
	case "03":
		return weather.IconNightScatteredClouds, nil
	case "04":
		return weather.IconNightBrokenClouds, nil
	case "09":
		return weather.IconNightShowerRain, nil
	case "10":
		return weather.IconNightRain, nil
	case "11":
		return weather.IconNightThunderstorm, nil
	case "13":
		return weather.IconNightSnow, nil
	case "50":
		return weather.IconNightMist, nil
	}
	return weather.IconNightClearSky, fmt.Errorf("nightIcon(): unknown icon %s", icon)
}
