package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
	"github.com/maximekuhn/metego/internal/calendar/sqlite"
	"github.com/maximekuhn/metego/internal/server"
	openweatherapi "github.com/maximekuhn/metego/internal/weather/open_weather_api"
	yaml "gopkg.in/yaml.v3"
)

// TODO:
// - add db file path in configFile
// - maybe add api key in the config
// - use a flag to specify the config file path

type configFile struct {
	Cities []string `yaml:"cities"`
}

func readConfig() (*configFile, error) {
	// TODO: read from a flag
	const (
		configFilePath string = "./config.metego.yaml"
	)

	buf, err := os.ReadFile(configFilePath)
	if err != nil {
		return nil, err
	}

	var cf configFile
	if err := yaml.Unmarshal(buf, &cf); err != nil {
		return nil, err
	}

	return &cf, nil
}

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

	conf, err := readConfig()
	var cities []string
	if err != nil {
		fmt.Printf("could not read config file: '%s'\n", err)
		fmt.Println("program will continue...")
		cities = make([]string, 0)
	} else {
		fmt.Println("successfully parsed config file")
		cities = conf.Cities
	}
	fmt.Printf("cities: %v\n", cities)

	// TODO: read from conf
	dbFilePath := "./metego.sqlite3"
	db, err := sql.Open("sqlite3", dbFilePath)
	if err != nil {
		panic(err)
	}
	defer func() { _ = db.Close() }()

	err = sqlite.ApplyMigrations(db)
	if err != nil {
		fmt.Println("failed to apply migrations")
		return
	}

	bdaysStorage, err := sqlite.NewSQLiteBirthdayStorage(db)
	if err != nil {
		fmt.Println("failed to create birthday storage")
		return
	}

	aptsStorage := sqlite.NewSQLiteAppointmentStorage(db)
	namedaysStorage := sqlite.NewSQLiteNamedayStorage(db)

	server := server.NewServer(fetcher, bdaysStorage, aptsStorage, namedaysStorage, cities)
	if err := server.Start(); err != nil {
		fmt.Println(err)
	}
}
