package main

import (
	"database/sql"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	_ "github.com/lib/pq" // Import the PostgreSQL driver

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"

	"github.com/team-scaletech/common/database"
	"github.com/team-scaletech/common/logging"

	_const "github.com/team-scaletech/common/const"
)

func main() {
	action := os.Args[1]
	switch action {
	case "migration":
		migrationFunc()
	case "migrate":
		migrateFunc()
	case "seeder":
		migrateFunc()
	default:
		panic("invalid action")
	}
}

func migrationFunc() {
	log := logging.GetLog()

	flag.Parse()
	fileName := os.Args[2]

	gooseCmd := exec.Command("goose", "-dir", "data_model/sql", "create", fileName, "sql")
	gooseCmd.Stdout = os.Stdout
	gooseCmd.Stderr = os.Stderr

	err := gooseCmd.Run()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create migration file")
	}

	fmt.Println(_const.Green, "Migration file created successfully")
}

func migrateFunc() {
	log := logging.GetLog()

	_, b, _, _ := runtime.Caller(0)
	basePath := filepath.Dir(b)

	bashPathFolder := strings.Split(basePath, "/")
	workspacePath := strings.Join(bashPathFolder[0:len(bashPathFolder)-1], "/")

	// Load the environment vars from a .env file if present
	// Get platform env path for database connection
	err := godotenv.Load(workspacePath + "/project/.env")
	if err != nil {
		log.Error().Err(err).Msg(err.Error())
	}

	type Config struct {
		Db database.DBConfig `required:"true"`
	}

	// Load the config struct with values from the environment without any prefix (i.e. "")
	var config Config
	err = envconfig.Process("", &config)
	if err != nil {
		log.Error().Err(err).Msg(err.Error())
	}

	connectionString := fmt.Sprintf("postgresql://%s:%s@/%s?host=%s&port=%d&sslmode=disable", config.Db.User, config.Db.Password, config.Db.Name, config.Db.Host, config.Db.Port)

	action := os.Args[2]
	switch action {
	case "up":
		log.Info().Msg("Upgrade migration version.!")
		err := executeGooseCommand(connectionString, "up")
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to apply migrations.")
		}
	case "down":
		log.Info().Msg("Downgrade migration version.!")
		err := executeGooseCommand(connectionString, "down")
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to revert migrations.")
		}
	case "force":
		log.Info().Msg("Force migration version.!")
		version := os.Args[3]
		err := executeGooseCommand("up", version)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to force migration version.")
		}
	case "seed":
		log.Info().Msg("Seed data!")
		seederFunc(connectionString)
	default:
		panic("Please select a valid action")
	}
}

func seederFunc(connectionString string) {
	log := logging.GetLog()

	_, b, _, _ := runtime.Caller(0)
	basePath := filepath.Dir(b)
	basePath = filepath.Join(basePath, "seeder")
	files, err := os.ReadDir(basePath)
	if err != nil {
		log.Error().Err(err).Msg(err.Error())
	}

	// Create a database connection
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Error().Err(err).Msg(err.Error())
	}
	defer db.Close()

	// Iterate through seed data files
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		seedFileName := file.Name()
		seedFilePath := filepath.Join(basePath, seedFileName)
		seedData, err := os.ReadFile(seedFilePath)
		if err != nil {
			log.Info().Msgf("Failed to read seed file %s: %v", seedFileName, err)
			continue
		}

		// Split seed data into individual insert statements
		insertStatements := strings.Split(string(seedData), ";")

		// Execute each insert statement
		for _, statement := range insertStatements {
			statement = strings.TrimSpace(statement)
			if statement == "" {
				continue
			}

			_, err = db.Exec(statement)
			if err != nil {
				log.Info().Msgf("Failed to insert seed data from file %s: %v", seedFileName, err)
				continue
			}
		}

		log.Info().Msgf("Seed data from file %s inserted successfully.", seedFileName)
	}
}

func executeGooseCommand(connectionString string, args ...string) error {
	gooseCmd := append([]string{"-allow-missing", "-dir", "data_model/sql", "postgres", connectionString}, args...)
	cmd := exec.Command("goose", gooseCmd...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
