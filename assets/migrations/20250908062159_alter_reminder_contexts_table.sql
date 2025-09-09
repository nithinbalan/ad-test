-- +goose Up
-- +goose StatementBegin
ALTER TABLE reminder_contexts
ADD module_id INT NULL;
ALTER TABLE reminder_contexts
ADD CONSTRAINT FK_reminder_contexts_modules
FOREIGN KEY (module_id)
REFERENCES modules(id)
ON DELETE NO ACTION ON UPDATE NO ACTION;
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
ALTER TABLE reminder_contexts
DROP CONSTRAINT FK_reminder_contexts_modules;
ALTER TABLE reminder_contexts
DROP COLUMN module_id;
-- +goose StatementEnd