package v1Ctl

import (
	"github.com/team-scaletech/common/validator"

	v1Srv "github.com/team-scaletech/project/services/v1"
)

type ApiCtl struct {
	BatchLogsCtl *BatchLogsCtl // Controller for handling batch logging operations
}

// InitV1BatchLogsCtl Controller
func InitV1BatchLogsCtl(validatorService validator.IAPIValidatorService, BatchLogsService v1Srv.IBatchLogsService) *BatchLogsCtl {
	BatchLogsCtl := BatchLogsCtl{
		BatchLogsService: BatchLogsService,
		APIValidator:     validatorService,
	}

	return &BatchLogsCtl
}
