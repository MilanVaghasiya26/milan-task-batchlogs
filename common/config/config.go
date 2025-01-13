package config

import (
	"time"

	"github.com/team-scaletech/common/database"
)

type JWTConfig struct {
	SecretKey        string `required:"true" split_words:"true"`
	ExpiryTimeInHour string `required:"true" split_words:"true"`
}

// Config is a struct representing the configuration for the application
type Config struct {
	Env         string            `required:"true"`  // Environment (e.g., "development", "production")
	Version     string            `required:"true"`  // Version of the application
	ServiceName string            `required:"true"`  // Name of the service
	ServicePort string            `required:"false"` // Port on which the service will run
	Path        string            `required:"false"` // Path information (optional)
	Level       string            `required:"false"` // Log level (optional)
	MaxAge      time.Duration     `required:"false"` // Maximum age (duration) for certain operations (optional)
	DB          database.DBConfig `required:"true"`  // Database configuration

	JWT JWTConfig `required:"true"`
}
