package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	AppPort         string
	DBHost          string
	DBPort          int
	DBUser          string
	DBPassword      string
	DBName          string
	DBDebug         bool
	BaseURLPath     string
	DBSSLMode       string
	RedisHost       string
	RedisPort       int
	RedisPassword   string
	RedisDB         int
	CartServicePort string
}

func LoadConfig() (*Config, error) {
	viper.New()
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("cannot read config file: %w", err)
	}

	config := &Config{
		AppPort:         viper.GetString("APP_PORT"),
		BaseURLPath:     viper.GetString("BASE_URL_PATH"),
		DBSSLMode:       viper.GetString("DB_SSL_MODE"),
		DBUser:          viper.GetString("DB_USER"),
		DBHost:          viper.GetString("DB_HOST"),
		DBPassword:      viper.GetString("DB_PASSWORD"),
		DBName:          viper.GetString("DB_NAME"),
		DBDebug:         viper.GetBool("DB_DEBUG"),
		DBPort:          viper.GetInt("DB_PORT"),
		RedisHost:       viper.GetString("REDIS_HOST"),
		RedisPort:       viper.GetInt("REDIS_PORT"),
		RedisPassword:   viper.GetString("REDIS_PASSWORD"),
		RedisDB:         viper.GetInt("REDIS_DB"),
		CartServicePort: viper.GetString("CART_SERVICE_PORT"),
	}

	return config, nil
}

func WriteTimeout() time.Duration {
	return 10 * time.Second
}

func ReadTimeout() time.Duration {
	return 10 * time.Second
}
