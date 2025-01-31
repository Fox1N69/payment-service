package config

import (
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

type Config struct {
	Env    EnvConfig    `mapstructure:"ENV_MODE"`
	Server ServerConfig `mapstructure:"server"`
	Psql   PSQLConfig   `mapstructure:"postgres"`
}

type EnvConfig struct {
	Mode string `mapstructure:"MODE"`
}

type ServerConfig struct {
	Port string `mapstructure:"port"`
}

type PSQLConfig struct {
	Host     string `mapstructure:"psql_host"`
	Port     string `mapstructure:"psql_port"`
	User     string `mapstructure:"psql_user"`
	Password string `mapstructure:"psql_password"`
	DBName   string `mapstructure:"psql_dbname"`
	SSLMode  string `mapstructure:"psql_sslmode"`
}

var (
	vprOnce sync.Once
	cfg     *Config
)

func LoadConfig(configFile string) *Config {
	vprOnce.Do(func() {
		if configFile != "" {
			err := godotenv.Load(configFile)
			if err != nil {
				log.Fatalf("[infra][Config][godotenv.Load] Error loading .env file: %v", err)
			}
		} else {
			err := godotenv.Load(".env")
			if err != nil {
				log.Fatalf("[infra][Config][godotenv.Load] Error loading .env file: %v", err)
			}
		}

		var config Config
		config.Env.Mode = getEnv("MODE", "development")
		config.Server.Port = getEnv("SERVER_PORT", "4000")
		config.Psql.Host = getEnv("PSQL_HOST", "localhost")
		config.Psql.Port = getEnv("PSQL_PORT", "5432")
		config.Psql.User = getEnv("PSQL_USER", "postgres")
		config.Psql.Password = getEnv("PSQL_PASSWORD", "")
		config.Psql.DBName = getEnv("PSQL_DBNAME", "")
		config.Psql.SSLMode = getEnv("PSQL_SSLMODE", "disable")

		cfg = &config
	})

	return cfg
}

func getEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}
