// Package config provides functionality to load and manage the application configuration.
package config

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"

	"user_segmentation_service/internal/db"
	"user_segmentation_service/internal/logger"
)

// Config holds the entire application configuration.
type Config struct {
	Log logger.Config `envconfig:"LOG" required:"true"`
	DB  db.Config     `envconfig:"DB" required:"true"`
}

// MustLoad is a function that loads environment variables from a `.env` file and
// populates a `Config` struct with the corresponding values. If the `.env` file
// cannot be loaded, the function will panic. It uses the `godotenv` package to
// load the environment variables and the `envconfig` package to process and
// map them to the `Config` struct. The function returns a pointer to the
// populated `Config` struct.
//
// Note: This function is designed to be used in scenarios where the application
// cannot proceed without the environment variables, hence the use of `panic`
// in case of failure.
func MustLoad() *Config {
	if err := godotenv.Load(".env"); err != nil {
		panic("cannot load .env file")
	}

	var config Config
	envconfig.MustProcess("", &config)

	return &config
}
