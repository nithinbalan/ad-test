-- +goose Up
-- +goose StatementBegin
ALTER TABLE incident_root_cause_analyses ALTER COLUMN question NVARCHAR(1000) NOT NULL;
ALTER TABLE incident_root_cause_analyses ALTER COLUMN answer NVARCHAR(MAX) NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE incident_root_cause_analyses ALTER COLUMN question NVARCHAR(500) NOT NULL;
ALTER TABLE incident_root_cause_analyses ALTER COLUMN answer NVARCHAR(1000) NULL;
-- +goose StatementEnd
