-- +goose Up
-- +goose StatementBegin

ALTER TABLE system_preferences
ADD theme_id INT NULL;

ALTER TABLE system_preferences
ADD CONSTRAINT FK_system_preferences_themes
FOREIGN KEY (theme_id)
REFERENCES themes (id)
ON DELETE NO ACTION
ON UPDATE NO ACTION;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE system_preferences
DROP CONSTRAINT FK_system_preferences_themes;

ALTER TABLE system_preferences
DROP COLUMN theme_id;

-- +goose StatementEnd