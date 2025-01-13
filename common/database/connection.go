package database

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/avast/retry-go/v3"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DBConfig is a struct representing the configuration parameters for the database connection
type DBConfig struct {
	User       string `required:"true" split_words:"true"`
	Password   string `required:"true" split_words:"true"`
	Host       string `required:"true" split_words:"true"`
	Port       int    `required:"true" split_words:"true"`
	Name       string `required:"true" split_words:"true"`
	disableSSL bool
}

// RowsAndCount is a struct representing the result of a database query with row count, ID, and error
type RowsAndCount struct {
	c   int
	id  *uuid.UUID
	err error
}

// pool is a global variable representing the connection pool for pgx
var (
	pool  *pgxpool.Pool
	poolG *gorm.DB
	once  sync.Once
)

// URL generates a connection URL using the DBConfig parameters
func (config *DBConfig) URL() string {
	url := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d",
		config.Host,
		config.User,
		config.Password,
		config.Name,
		config.Port,
	)

	if config.disableSSL {
		url = fmt.Sprintf("%s sslmode=disable", url)
	}

	return url
}

// getConfig generates a pgxpool.Config from the provided DBConfig
func getConfig(dbConfig DBConfig) (*pgxpool.Config, error) {
	config, err := pgxpool.ParseConfig(dbConfig.URL())
	if err != nil {
		return nil, err
	}

	config.ConnConfig.LogLevel = pgx.LogLevelDebug

	return config, nil
}

// resetPool resets the global connection pool
func resetPool() {
	pool = nil
	once = sync.Once{}
}

// CreateConnectionPool creates a new connection pool based on the provided DBConfig
func CreateConnectionPool(dbConfig DBConfig) (*pgxpool.Pool, error) {
	config, err := getConfig(dbConfig)
	if err != nil {
		return nil, err
	}
	// Set the maximum number of connections in the pool
	config.MaxConns = int32(10)
	return createConnectionPool(config)
}

// createConnectionPool creates a connection pool with the given pgxpool.Config
func createConnectionPool(config *pgxpool.Config) (*pgxpool.Pool, error) {
	var err error
	once.Do(func() {
		err = retry.Do(func() error {
			p, err := pgxpool.ConnectConfig(context.Background(), config)
			if err != nil {
				return err
			}
			pool = p
			return nil
		},
			retry.Delay(1*time.Second))
		if err != nil {
			log.Fatal(err.Error())
		}
	})

	return pool, err
}

// GetConnectionPool returns the global connection pool
func GetConnectionPool() (*pgxpool.Pool, error) {
	if pool == nil {
		return nil, errors.New("connection Pool has not been created")
	}
	return pool, nil
}

// db is a global variable representing the GORM database instance
var db *gorm.DB
var connectionOnce sync.Once

// Init initializes the GORM database connection
func Init(config *DBConfig) {
	connectionOnce.Do(func() {
		var err error
		connectionString := config.URL()
		logMode := logger.LogLevel(4)
		db, err = gorm.Open(postgres.Open(connectionString), &gorm.Config{
			Logger: logger.Default.LogMode(logMode),
		})
		if err != nil {
			panic(err)
		}
		poolG = db
	})
}

// GetDB returns the global GORM database instance
func GetDB() *gorm.DB {
	return db
}
