-- +goose Up
-- +goose StatementBegin
DROP TABLE IF EXISTS logs;

CREATE TABLE logs (
    id UUID DEFAULT uuid_generate_v1mc(),
    timestamp TIMESTAMP WITH TIME ZONE NOT NULL,
    body TEXT NOT NULL,
    service CHARACTER VARYING(255) NOT NULL,
    severity CHARACTER VARYING(50) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    CONSTRAINT logs_pkey PRIMARY KEY (id)
);

COMMENT ON COLUMN logs.id IS 'Unique identifier for the log entry';

COMMENT ON COLUMN logs.timestamp IS 'Timestamp of the log entry';

COMMENT ON COLUMN logs.body IS 'Content of the log entry';

COMMENT ON COLUMN logs.service IS 'The service associated with the log entry';

COMMENT ON COLUMN logs.severity IS 'The severity level of the log entry (e.g., info, warning, error)';

COMMENT ON COLUMN logs.created_at IS 'Timestamp when the log entry was created';

COMMENT ON COLUMN logs.updated_at IS 'Timestamp when the log entry was last updated';

COMMENT ON COLUMN logs.deleted_at IS 'Timestamp when the log entry was deleted or marked as deleted';

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS logs;

-- +goose StatementEnd