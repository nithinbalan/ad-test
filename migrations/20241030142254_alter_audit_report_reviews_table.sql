-- +goose Up
-- +goose StatementBegin
ALTER TABLE audit_report_reviews
ALTER COLUMN audit_draft_report_id INT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE audit_report_reviews
ALTER COLUMN audit_draft_report_id INT NOT NULL;
-- +goose StatementEnd
