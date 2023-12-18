package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

const version = "BUILD_VERSION"

type App struct {
	Addr string
	Port int
}

type Postgres struct {
	DSN string
}

type Nats struct {
	ClusterID string
	ClientID  string
	URL       string
}

type Config struct {
	App      App
	Postgres Postgres
	Nats     Nats
}

// New returns a new Config struct
func New() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	return &Config{
		App: App{
			Addr: getEnv("APP_ADDR"),
			Port: getEnvAsInt("APP_PORT"),
		},
		Postgres: Postgres{
			DSN: getEnv("POSTGRES_DSN"),
		},
		Nats: Nats{
			ClusterID: getEnv("NATS_CLUSTER_ID"),
			ClientID:  getEnv("NATS_CLIENT_ID"),
			URL:       getEnv("NATS_URL"),
		},
	}, nil
}

func getEnv(key string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return ""
}

func getEnvAsInt(name string) int {
	valueStr := getEnv(name)
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return 0
}
