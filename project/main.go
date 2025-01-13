package main

import (
	"context"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/team-scaletech/common/config"
	"github.com/team-scaletech/common/database"
	"github.com/team-scaletech/common/helpers"
	"github.com/team-scaletech/common/logging"
	"github.com/team-scaletech/project/routers/api"

	constants "github.com/team-scaletech/project/utils/const"
)

// Constants for default values
const (
	defaultPort = constants.DefaultPort
)

// Execution starts from main function
func main() {
	zlog := logging.GetLog()

	// Load the config struct with values from the environment without any prefix (i.e. "")
	var cf config.Config

	// Load the environment vars from a .env file if present
	err := godotenv.Load()
	if err != nil {
		zlog.Error().Err(err).Msg("")
	}

	// Process environment variables
	err = envconfig.Process("", &cf)
	if err != nil {
		zlog.Error().Err(err).Msg("")
		panic("Failed to get env: " + err.Error())
	}

	// Get the port from the environment variable.
	port := cf.ServicePort

	// If the `port` is not specified in the configuration, use the `defaultPort` as a fallback.
	if port == "" {
		cf.ServicePort = defaultPort
	}

	// The database connection pool must be initialized before it can be used
	database.Init(&cf.DB)
	_, err = database.CreateConnectionPool(cf.DB)
	if err != nil {
		panic(err.Error())
	}
	zlog.Info().Msg("Database connected successfully.")

	// Create a new api instance with the provided configuration.
	rt := api.NewRouter(cf)

	// Set up the router, initializing routes and other configurations.
	rt.Setup()

	// Start the router in a separate goroutine.
	go rt.Run()

	// Gracefully stop the router and perform necessary cleanup.
	helpers.GracefulStop(func(ctx context.Context) error {
		var err error
		// Close the router, handling any potential errors.
		if err = rt.Close(ctx); err != nil {
			return err
		}
		zlog.Info().Msg("Server shutdown gracefully.")
		return nil
	})
}
