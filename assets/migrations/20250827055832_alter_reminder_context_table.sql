-- +goose Up
-- +goose StatementBegin
ALTER TABLE reminder_contexts
ADD status VARCHAR(10) NOT NULL DEFAULT 'active',
    CONSTRAINT CHK_reminder_contexts_status CHECK (status IN ('active', 'inactive'));
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE reminder_contexts
DROP CONSTRAINT CHK_reminder_contexts_status;

ALTER TABLE reminder_contexts
DROP COLUMN status;
-- +goose StatementEnd
