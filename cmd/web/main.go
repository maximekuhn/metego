package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/maximekuhn/metego/server"
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
	fetcher := openweatherapi.NewOpenWeatherFetcher(apiKey)

	server := server.NewServer(fetcher)
	if err := server.Start(); err != nil {
		fmt.Println(err)
	}
}
