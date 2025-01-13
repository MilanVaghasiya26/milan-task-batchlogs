package v1Ctl

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/team-scaletech/common/helpers"
	"github.com/team-scaletech/common/logging"
	"github.com/team-scaletech/project/utils/message"

	validator "github.com/team-scaletech/common/validator"
	v1Req "github.com/team-scaletech/project/resources/request/v1"
	v1Srv "github.com/team-scaletech/project/services/v1"
)

type BatchLogsCtl struct {
	BatchLogsService v1Srv.IBatchLogsService        // Service interface for batch logs operations.
	APIValidator     validator.IAPIValidatorService // Validator service to validate API requests.
}

// @Summary		Ingest logs
// @Description	API to ingest logs into the system.
// @Tags			BatchLogs
// @Accept			json
// @Produce		json
// @Param			request	body		v1Req.LogsEntry				true	"Logs Entry Payload"
// @Success		201		{object}	helpers.ResponseEntities	"Batch logs stored successfully"
// @Failure		400		{object}	helpers.ErrorResponse		"Bad request or validation error"
// @Failure		500		{object}	helpers.ErrorResponse		"Internal server error"
// @Router			/platform/api/v1/ingest [post]
func (blc *BatchLogsCtl) CreateBatchLogs(c *gin.Context) {
	// Retrieve the request-specific logger for debugging.
	log := logging.GetRequestLog(c)

	// Parse the request body into the LogsEntry struct.
	var req v1Req.LogsEntry

	// Decode the JSON request body into the LogsEntry struct.
	if err := c.BindJSON(&req); err != nil {
		msg := fmt.Sprintf("error decoding expected request parameters: %s", err)
		log.Error().Err(err).Msg(msg)
		helpers.BadRequest(c, helpers.ErrorResponse{Message: msg})
		return
	}

	// Validate the parsed LogsEntry struct.
	if _, ok := blc.APIValidator.ValidateStruct(req); !ok {
		log.Error().Msgf("Error Validate struct: %s", helpers.PrettyPrinter(c, req))
		helpers.BadRequest(c, helpers.ErrorResponse{Message: message.SomethingWrong})
		return
	}

	// Call the BatchLogsCreate service method to persist the log entry.
	err := blc.BatchLogsService.BatchLogsCreate(c, req)
	if err != nil {
		helpers.ServiceErrorResponse(c, err)
	}

	// Respond with a 201 Created status and success message.
	helpers.StatusCreated(c, &helpers.ResponseEntities{Message: message.CreateBatchLogsSuccess})
}

// @Summary		Query logs
// @Description	API to query logs based on time range and search text.
// @Tags			BatchLogs
// @Accept			json
// @Produce		json
// @Param			start	query		string						false	"Start date or timestamp (optional)"
// @Param			end		query		string						false	"End date or timestamp (optional)"
// @Param			text	query		string						false	"Search text to filter logs (optional)"
// @Success		200		{object}	helpers.ResponseEntities	"Batch logs fetched successfully"
// @Failure		400		{object}	helpers.ErrorResponse		"Bad request"
// @Failure		500		{object}	helpers.ErrorResponse		"Internal server error"
// @Router			/platform/api/v1/query [get]
func (blc *BatchLogsCtl) ListBatchLogs(c *gin.Context) {
	// Extract query parameters for filtering logs.
	start := c.Query("start")     // Start date or timestamp.
	end := c.Query("end")         // End date or timestamp.
	searchText := c.Query("text") // Text to filter logs by.

	// Clean up whitespace in query parameters.
	start = strings.TrimSpace(start)
	end = strings.TrimSpace(end)
	searchText = strings.TrimSpace(searchText)

	// Call the BatchLogsList service method to fetch the logs based on filters.
	err := blc.BatchLogsService.BatchLogsList(c, start, end, searchText)
	if err != nil {
		helpers.ServiceErrorResponse(c, err)
	}

	// Respond with a 200 OK status and success message.
	helpers.StatusOK(c, &helpers.ResponseEntities{Message: message.ListBatchLogsSuccess})
}
