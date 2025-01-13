package v1Repo

import (
	"database/sql"

	"github.com/team-scaletech/common/database"
	"github.com/team-scaletech/data_model/model"
)

type IBatchLogsRepository interface {
	BatchLogsCreate(req *model.LogEntry) error
}

type BatchLogsRepo struct {
	DB *sql.DB
}

func NewBatchLogsWriter() IBatchLogsRepository {
	return &BatchLogsRepo{}
}

// BatchLogsCreate inserts a log entry into the database.
func (ur *BatchLogsRepo) BatchLogsCreate(logEntry *model.LogEntry) error {
	// Use the ORM to insert the log entry and return any errors.
	return database.GetDB().Create(logEntry).Error
}
