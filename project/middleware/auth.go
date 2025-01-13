package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"

	"github.com/team-scaletech/common/config"

	"github.com/team-scaletech/common/helpers"
)

// IMiddleware is an interface defining the methods that middleware should implement
type IMiddleware interface {
	AuthHandler(jwtSecret string) gin.HandlerFunc
}

// Middleware is a struct representing middleware configuration
type Middleware struct {
	Config config.Config
}

// NewUserMiddlewareService creates a new instance of Middleware with the provided configuration
func NewUserMiddlewareService(cf config.Config) IMiddleware {
	return &Middleware{
		Config: cf,
	}
}

// AuthHandler is a middleware function that handles user authentication
func (m *Middleware) AuthHandler(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract token from the "Authorization" header
		token := c.Request.Header.Get("Authorization")
		b := "Bearer "
		if !strings.Contains(token, b) {
			// If the token doesn't contain "Bearer ", respond with an unauthorized message
			helpers.StatusUnauthorized(c, helpers.ErrorResponse{Message: "Your request is not authorized"})
			c.Abort()
			return
		}

		// Split the token to get the actual token value
		t := strings.Split(token, b)
		if len(t) < 2 {
			// If there is no token value, respond with an unauthorized message
			helpers.StatusUnauthorized(c, helpers.ErrorResponse{Message: "An authorization token was not supplied"})
			c.Abort()
			return
		}

		// Validate the token
		valid, err := ValidateToken(t[1], jwtSecret)
		if err != nil {
			// If the token is invalid, respond with an unauthorized message
			helpers.StatusUnauthorized(c, helpers.ErrorResponse{Message: "Invalid authorization token"})
			c.Abort()
			return
		}

		// Set "userData" variable in the Gin context
		c.Set("userData", valid.Claims.(jwt.MapClaims)["userData"])
		c.Next()
	}
}

// ValidateToken validates a JWT token using the provided key
func ValidateToken(t string, k string) (*jwt.Token, error) {
	token, err := jwt.Parse(t, func(token *jwt.Token) (interface{}, error) {
		return []byte(k), nil
	})

	return token, err
}
