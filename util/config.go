package util

import (
	"time"

	"github.com/spf13/viper"
)

// stores all configuration of the application from file or env variables
type Config struct {
	DBDriver           string        `mapstructure:"DB_DRIVER"`
	DBSource           string        `mapstructure:"DB_SOURCE"`
	ServerAddress      string        `mapstructure:"SERVER_ADDRESS"`
	JWTSigningKey      string        `mapstructure:"JWT_SIGNING_KEY"`
	PasetoSymmetricKey string        `mapstructure:"PASETO_SYMMETRIC_KEY"`
	TokenDuration      time.Duration `mapstructure:"TOKEN_DURATION"`
}

// reads configuration files and returns the config object
func LoadConfig(path string) (Config, error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	//looks for the env files that match the above
	viper.AutomaticEnv()

	//reads the config and loads
	err := viper.ReadInConfig()
	if err != nil {
		return Config{}, err
	}

	//unmarshall the configs into the config object
	var config Config
	err = viper.Unmarshal(&config)
	if err != nil {
		return Config{}, err
	}

	return config, nil
}
