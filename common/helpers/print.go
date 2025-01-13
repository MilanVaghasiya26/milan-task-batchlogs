package helpers

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/team-scaletech/common/logging"
)

func PrettyPrinter(c *gin.Context, data interface{}) string {
	log := logging.GetRequestLog(c)

	b, err := json.Marshal(data)
	if err != nil {
		log.Error().Err(err).Msgf("Unable to Marshal data")
		return ""
	}
	return string(b)
}
