package main

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	openweatherapi "github.com/maximekuhn/metego/weather/open_weather_api"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	apiKey := os.Getenv("OPEN_WEATHER_API_KEY")
	if len(apiKey) == 0 {
		panic("OPEN_WEATHER_API_KEY is empty")
	}

	city := os.Getenv("CITY")
	if len(city) == 0 {
		panic("CITY is empty")
	}

	fetcher := openweatherapi.NewOpenWeatherFetcher(apiKey)

	for {
		curr, err := fetcher.FetchCurrent(city)
		if err != nil {
			fmt.Printf("failed to fetch current weather for %s: %s\n", city, err)
			continue
		}

		fmt.Printf("Current temperature for %s: %.2f°C\n", city,curr.Temp)

		time.Sleep(10 * time.Second)
	}
}
