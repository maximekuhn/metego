package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
	"github.com/maximekuhn/metego/calendar/sqlite"
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

	// TODO: read from conf
	dbFilePath := "./metego.sqlite3"
	db, err := sql.Open("sqlite3", dbFilePath)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = sqlite.ApplyMigrations(db)
	if err != nil {
		fmt.Println("failed to apply migrations")
		return
	}

	storage, err := sqlite.NewSQliteBdayStorage(db)
	if err != nil {
		fmt.Println("failed to create birthday storage")
		return
	}

	server := server.NewServer(fetcher, storage)
	if err := server.Start(); err != nil {
		fmt.Println(err)
	}
}
