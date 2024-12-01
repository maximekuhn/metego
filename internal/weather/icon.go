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

	IconNightClearSky
	IconNightFewClouds
	IconNightScatteredClouds
	IconNightBrokenClouds
	IconNightShowerRain
	IconNightRain
	IconNightThunderstorm
	IconNightSnow
	IconNightMist
)
