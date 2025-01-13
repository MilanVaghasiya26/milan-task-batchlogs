package crypto

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"

	"github.com/team-scaletech/common/config"
	"github.com/team-scaletech/common/helpers"
	"github.com/team-scaletech/common/logging"
)

// tokenData is a struct to hold the data that will be used to generate an authentication token
type tokenData struct {
	Id        uuid.UUID `json:"id"`         // UUID is a string representation of the unique identifier
	CreatedAt time.Time `json:"created_at"` // CreatedAt is the timestamp indicating when the token was created
}

// GenerateAuthToken generates an authentication token using the provided UUID and creation timestamp
func GenerateAuthToken(id uuid.UUID, createdAt time.Time, config config.Config) string {
	zlog := logging.GetLog()

	// Create a tokenData instance with the UUID and creation timestamp
	tokenData := &tokenData{
		Id:        id,
		CreatedAt: createdAt,
	}

	jwtExpiry, err := strconv.Atoi(config.JWT.ExpiryTimeInHour)
	if err != nil {
		zlog.Error().Err(err).Msgf("Error while converting string value of jwt expiry to int:%+v", err.Error())
	}

	// Generate a token using the middlewares GenerateToken function
	token, err := GenerateToken([]byte(config.JWT.SecretKey), tokenData, jwtExpiry)
	if err != nil {
		// Log a warning message if there is an error generating the token
		zlog.Warn().Err(err).Msg(err.Error())
	}

	// Return the generated token
	return token
}

// GenerateToken creates a new JWT token with the provided key and user data
func GenerateToken(k []byte, userData interface{}, jwtExpiry int) (string, error) {
	// Create the token
	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims including user data and expiration time
	claims := make(jwt.MapClaims)
	claims["userData"] = userData
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(jwtExpiry)).Unix()
	token.Claims = claims

	// Sign and get the complete encoded token as a string
	tokenString, err := token.SignedString(k)

	return tokenString, err
}

func GetIdFromToken(c *gin.Context) (uuid.UUID, error) {
	log := logging.GetRequestLog(c)

	var userUUID uuid.UUID

	userInterface, exists := c.Get("userData")
	if !exists {
		message := "authenticated user with user_id required"
		log.Warn().Msgf("%v: exists: %v", message, exists)

		return userUUID, fmt.Errorf(message)
	}

	var response tokenData

	bytesData, err := json.Marshal(userInterface)
	if err != nil {
		// Handle the error if the conversion fails
		log.Error().Err(err).Msgf("Error while marshaling user token data: %s", err)
		helpers.StatusInternalServerError(c, helpers.ErrorResponse{Message: err.Error()})
		return userUUID, err
	}
	err = json.Unmarshal(bytesData, &response)
	if err != nil {
		log.Error().Err(err).Msgf("Error while unmarshalling user token bytes data: %s", err)
		helpers.StatusInternalServerError(c, helpers.ErrorResponse{Message: err.Error()})
		return userUUID, err
	}

	userId := response.Id
	return userId, nil
}
