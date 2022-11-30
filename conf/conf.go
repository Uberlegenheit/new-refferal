package conf

import (
	"os"
	"strconv"
	"strings"
)

type (
	Config struct {
		API      API
		LogLevel string
		Postgres Postgres
		Cloud    Cloud
	}
	API struct {
		ListenOnPort       uint64
		CORSAllowedOrigins []string
		HttpSwaggerAddress string
	}
	Postgres struct {
		Host     string
		User     string
		Password string
		Database string
		SSLMode  string
	}
	Cloud struct {
		Name   string
		Key    string
		Secret string
	}
)

const (
	Service = "hr-board"
)

func GetNewConfig() (cfg Config, err error) {
	port, err := strconv.ParseInt(os.Getenv("LISTEN_PORT"), 10, 64)
	if err != nil {
		return Config{}, err
	}
	return Config{
		API: API{
			ListenOnPort:       uint64(port),
			CORSAllowedOrigins: strings.Split(os.Getenv("CORS_ALLOWED"), ","),
		},
		LogLevel: os.Getenv("LOG_LEVEL"),
		Postgres: Postgres{
			Host:     os.Getenv("DB_HOST"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			Database: os.Getenv("DB_NAME"),
			SSLMode:  os.Getenv("DB_SSL"),
		},
		Cloud: Cloud{
			Name:   os.Getenv("CLOUD_NAME"),
			Key:    os.Getenv("CLOUD_KEY"),
			Secret: os.Getenv("CLOUD_SECRET"),
		},
	}, nil
}
