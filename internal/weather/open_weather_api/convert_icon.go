package openweatherapi

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/maximekuhn/metego/internal/weather"
)

// try to convert the icon string from open weather API to a custom type
//
// if an error si returned ,the icon is not valid
func toWeatherIcon(icon string) (weather.CurrentWeatherIcon, error) {
	// remove 'd' or 'n' to only keep icon id's number
	if len(icon) != 3 {
		return 0, errors.New("unknown icon type")
	}

	iconID, err := strconv.ParseInt(icon[:2], 10, 8)
	if err != nil {
		return 0, err
	}

	iconIDNumber := int8(iconID)
	switch iconIDNumber {
	case 1:
		return weather.IconClearSky, nil
	case 2:
		return weather.IconFewClouds, nil
	case 3:
		return weather.IconScatteredClouds, nil
	case 4:
		return weather.IconBrokenClouds, nil
	case 9:
		return weather.IconShowerRain, nil
	case 10:
		return weather.IconRain, nil
	case 11:
		return weather.IconThunderstorm, nil
	case 13:
		return weather.IconSnow, nil
	case 50:
		return weather.IconMist, nil
	}

	return 0, fmt.Errorf("unknown icon ID: %s", icon)
}
