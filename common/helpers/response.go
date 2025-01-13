package helpers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

func GetError(field, tag string) string {
	if tag == "required" {
		return field + " is " + tag
	} else if tag == "email" {
		return field + " not valid"
	}
	return ""
}

func BadRequest(c *gin.Context, errResp ErrorResponse) {
	c.JSON(http.StatusBadRequest, gin.H{
		"message": errResp.Message,
	})

}

func StatusUnprocessableEntity(c *gin.Context, errResp ErrorResponse) {
	c.JSON(http.StatusUnprocessableEntity, gin.H{
		"message": errResp.Message,
	})
}

func StatusUnauthorized(c *gin.Context, errResp ErrorResponse) {
	c.JSON(http.StatusUnauthorized, gin.H{
		"message": errResp.Message,
	})
}

func StatusNotFound(c *gin.Context, errResp ErrorResponse) {
	c.JSON(http.StatusNotFound, gin.H{
		"message": errResp.Message,
	})
}

func StatusForbidden(c *gin.Context, errResp ErrorResponse) {
	c.JSON(http.StatusForbidden, gin.H{
		"message": errResp.Message,
	})
}

func StatusConflict(c *gin.Context, errResp ErrorResponse) {
	c.JSON(http.StatusConflict, gin.H{
		"message": errResp.Message,
	})
}

func StatusInternalServerError(c *gin.Context, errResp ErrorResponse) {
	c.JSON(http.StatusInternalServerError, gin.H{
		"message": errResp.Message,
	})
}

func StatusNoContent(c *gin.Context, errResp ErrorResponse) {
	c.JSON(http.StatusNoContent, gin.H{
		"message": errResp.Message,
	})
}

func StatusOK(c *gin.Context, data *ResponseEntities) {
	c.JSON(http.StatusOK, data)
}

func StatusCreated(c *gin.Context, data *ResponseEntities) {
	c.JSON(http.StatusCreated, data)
}

func ServiceErrorResponse(c *gin.Context, err error) {
	var serviceError ServiceError
	if errors.As(err, &serviceError) {
		c.JSON(serviceError.Code, gin.H{
			"message": serviceError.Message,
		})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
	}
	c.Abort()
}
