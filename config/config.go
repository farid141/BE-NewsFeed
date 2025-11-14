package config

import "github.com/spf13/viper"

type Config struct {
	DBHost     string
	DBPort     int
	DBName     string
	DBUser     string
	DBPassword string
	LOG_FILE   string
}

func LoadConfig() (Config, error) {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		return Config{}, err
	}

	cfg := Config{
		DBHost:     viper.GetString("DB_HOST"),
		DBPort:     viper.GetInt("DB_PORT"),
		DBUser:     viper.GetString("DB_USER"),
		DBName:     viper.GetString("DB_NAME"),
		DBPassword: viper.GetString("DB_PASSWORD"),
		LOG_FILE:   viper.GetString("LOG_FILE"),
	}

	return cfg, nil
}
