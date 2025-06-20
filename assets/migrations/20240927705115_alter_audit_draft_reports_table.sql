-- +goose Up
-- +goose StatementBegin

-- Adding new field for next review level
ALTER TABLE audit_draft_reports
ADD next_review_level INT DEFAULT 0 NULL;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Dropping the next_review_level field
ALTER TABLE audit_draft_reports
DROP COLUMN next_review_level;

-- +goose StatementEnd
