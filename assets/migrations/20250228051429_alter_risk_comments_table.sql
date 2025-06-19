-- +goose Up
-- +goose StatementBegin
ALTER TABLE risk_comments ALTER COLUMN risk_timeline_type_id INT NULL;
ALTER TABLE risk_comments ALTER COLUMN risk_status_id INT NULL;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE risk_comments ALTER COLUMN risk_status_id INT NOT NULL;
ALTER TABLE risk_comments ALTER COLUMN risk_timeline_type_id INT NOT NULL;
-- +goose StatementEnd
