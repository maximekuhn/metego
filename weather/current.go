package weather

type CurrentWeatherIcon int

const (
	IconClearSky CurrentWeatherIcon = iota
	IconFewClouds
	IconScatteredClouds
	IconBrokenClouds
	IconShowerRain
	IconRain
	IconThunderstorm
	IconSnow
	IconMist
)

type CurrentWeather struct {
	Temp        float64
	FeelsLike   float64
	Pressure    float64
	Humidity    float64
	Description string
	Icon        CurrentWeatherIcon
}
