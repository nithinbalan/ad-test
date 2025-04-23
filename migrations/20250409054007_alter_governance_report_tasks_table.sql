-- +goose Up
-- +goose StatementBegin
ALTER TABLE governance_report_tasks
ALTER COLUMN target_date DATE NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE governance_report_tasks
ALTER COLUMN target_date DATE NOT NULL;
-- +goose StatementEnd
