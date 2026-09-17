-- +goose Up
-- +goose StatementBegin
ALTER TABLE incident_investigation_findings ALTER COLUMN title NVARCHAR(MAX) NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE incident_investigation_findings ALTER COLUMN title NVARCHAR(500) NOT NULL;
-- +goose StatementEnd
