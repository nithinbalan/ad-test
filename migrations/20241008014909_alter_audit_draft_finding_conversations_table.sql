-- +goose Up
-- +goose StatementBegin

ALTER TABLE audit_draft_finding_conversations
ADD is_accepted BIT NOT NULL DEFAULT 0;

ALTER TABLE audit_draft_finding_conversations
ADD is_dismissed BIT NOT NULL DEFAULT 0;

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin

ALTER TABLE audit_draft_finding_conversations
DROP COLUMN is_accepted;

ALTER TABLE audit_draft_finding_conversations
DROP COLUMN is_dismissed;

-- +goose StatementEnd
