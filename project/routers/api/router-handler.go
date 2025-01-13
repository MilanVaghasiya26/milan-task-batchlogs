package api

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/mvrilo/go-redoc"
	"github.com/thinkerou/favicon"

	ginredoc "github.com/mvrilo/go-redoc/gin"

	"github.com/team-scaletech/common/config"
	"github.com/team-scaletech/common/logging"
	"github.com/team-scaletech/common/validator"
	"github.com/team-scaletech/project/middleware"

	commonMiddleware "github.com/team-scaletech/common/middleware"
	v1Ctl "github.com/team-scaletech/project/controllers/v1"
	v1Srv "github.com/team-scaletech/project/services/v1"
)

// IRoutes is an interface that defines the methods required for setting up, running, and closing routes.
type IRoutes interface {
	Setup()
	Run()
	Close(ctx context.Context) error
}

// Routes is a struct that implements the IRoutes interface and holds the necessary components for handling routes.
type Routes struct {
	router     *gin.Engine            // Gin router instance
	server     *http.Server           // HTTP server instance
	config     config.Config          // Configuration settings
	middleware middleware.IMiddleware // Middleware interface for handling request/response processing
	apiCtl     *v1Ctl.ApiCtl
}

// NewRouter creates and initializes a new router instance with the provided configuration.
func NewRouter(config config.Config) IRoutes {
	// Initialize an API validator service for request validation.
	validation := validator.NewAPIValidatorService()

	// Initialize a user middleware service with the provided configuration.
	mw := middleware.NewUserMiddlewareService(config)

	// Initialize an api service/controller
	batchLogsSrv := v1Srv.NewBatchLogsService(config)
	batchLogsCtl := v1Ctl.InitV1BatchLogsCtl(validation, batchLogsSrv)

	apiControllers := &v1Ctl.ApiCtl{
		BatchLogsCtl: batchLogsCtl,
	}

	// Create a Gin router with default settings.
	router := gin.Default()

	// Create an HTTP server with the provided configuration.
	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", config.ServicePort),
		Handler: router,
	}

	// Return a Routes instance with the initialized components.
	return &Routes{
		router,
		server,
		config,
		mw,
		apiControllers,
	}
}

// Run starts the server and listens for incoming requests.
func (rt *Routes) Run() {
	zlog := logging.GetLog()
	zlog.Info().Msgf("Project service listen on %s", rt.config.ServicePort)
	err := rt.server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		zlog.Fatal().Err(err).Msgf("listen: %s\n", err)
	}
}

// Close gracefully shuts down the server and releases any allocated resources.
func (rt *Routes) Close(ctx context.Context) error {
	if rt.server != nil {
		return rt.server.Shutdown(ctx)
	}
	return nil
}

// Setup configures and sets up the application's routes, middleware and endpoints for the application.
func (rt *Routes) Setup() {
	zlog := logging.GetLog()

	// Add a request id to each incoming request
	rt.router.Use(commonMiddleware.DefaultRequestId())

	// Use logging middleware for request logging
	rt.router.Use(logging.Middleware)

	// Configure Cross-Origin Resource Sharing (CORS) settings.
	rt.setupCors()

	// Initialize and set up default endpoints based on the environment.
	rt.setupDefaultEndpoints()

	// Configure platform routes for version 1 (V1) of the application.
	PlatformRoutesV1(rt)

	basePath := fmt.Sprintf("/%s", rt.config.ServiceName)

	// Serve static HTML page
	// You should use gin's Static method to serve static files like index.html
	// rt.router.Static("/static", "./static") // serves files in ./static directory at /static URL path
	rt.router.GET("/", func(c *gin.Context) {
		c.File("./static/index.html") // serve the index.html file at the root
	})

	// Add the handler to serve the redoc
	specFile := "./docs/swagger.json"
	if _, err := os.Stat(specFile); err == nil {
		docs := redoc.Redoc{
			Title:       "Docs",
			Description: "Documentation",
			SpecFile:    specFile,
			SpecPath:    fmt.Sprintf("%s/docs/openapi.json", basePath),
			DocsPath:    fmt.Sprintf("%s/docs", basePath),
		}
		rt.router.Use(ginredoc.New(docs))
	} else {
		zlog.Warn().Msgf("Swagger file not found at %s, skipping redoc init", specFile)
	}

	// Serve a favicon to keep the logs clean
	rt.router.Use(favicon.New("./favicon.ico"))
}

// setupCors configures Cross-Origin Resource Sharing (CORS) settings for the router.
func (rt *Routes) setupCors() {
	rt.router.Use(cors.New(cors.Config{
		ExposeHeaders:   []string{"Data-Length"},
		AllowMethods:    []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:    []string{"Content-Type", "Authorization"},
		AllowAllOrigins: true,
		MaxAge:          12 * time.Hour,
	}))
}

// setupDefaultEndpoints sets up default endpoints for the Routes instance.
func (rt *Routes) setupDefaultEndpoints() {
	// Define a route for handling "/ping" endpoint
	rt.router.GET("/ping", func(c *gin.Context) {
		var msg string
		// Check if the environment is production
		if rt.config.Env == "production" {
			// Provide a detailed response for production environment
			msg = fmt.Sprintf("Pong! I am %s db for user project. Version is %s.", rt.config.Env, rt.config.Version)
		} else {
			// For non-production environments, a simple "pong" response
			msg = "pong"
		}
		c.JSON(200, gin.H{"message": msg})
	})

	// Define a route for handling "/health" endpoint
	rt.router.GET("/health", func(c *gin.Context) {
		// Respond with a plain text "OK" for health check
		c.Data(http.StatusOK, "text/plain; charset=utf-8", []byte("OK"))
	})
}
