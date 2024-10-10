-- +goose Up
-- +goose StatementBegin
ALTER TABLE audit_draft_finding_conversations
DROP COLUMN description;

ALTER TABLE audit_draft_finding_conversations
ADD comment NVARCHAR(MAX) NULL;
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
ALTER TABLE audit_draft_finding_conversations
DROP COLUMN comment;

ALTER TABLE audit_draft_finding_conversations
ADD description NVARCHAR(MAX) NULL;
-- +goose StatementEnd
