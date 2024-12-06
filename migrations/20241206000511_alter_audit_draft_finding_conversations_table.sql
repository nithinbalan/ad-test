-- +goose Up
-- +goose StatementBegin

ALTER TABLE audit_draft_finding_conversations
ADD is_reopened BIT NOT NULL DEFAULT 0;

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin

ALTER TABLE audit_draft_finding_conversations
DROP COLUMN is_reopened;

-- +goose StatementEnd
