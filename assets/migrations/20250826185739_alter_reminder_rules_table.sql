-- +goose Up
-- +goose StatementBegin
ALTER TABLE reminder_rules
ADD reminder_context_id INT NULL;

ALTER TABLE reminder_rules
ADD CONSTRAINT FK_reminder_rules_reminder_contexts
FOREIGN KEY (reminder_context_id)
REFERENCES reminder_contexts(id)
ON DELETE NO ACTION ON UPDATE NO ACTION;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE reminder_rules
DROP CONSTRAINT FK_reminder_rules_reminder_contexts;

ALTER TABLE reminder_rules
DROP COLUMN reminder_context_id;
-- +goose StatementEnd
