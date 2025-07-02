-- +goose Up
-- +goose StatementBegin

-- 1. Add the new column
ALTER TABLE active_directory_settings
ADD integration_id INT NULL;

-- 2. Add the foreign key constraint linking to integrations(id)
ALTER TABLE active_directory_settings
ADD CONSTRAINT FK_active_directory_settings_integrations
FOREIGN KEY (integration_id) REFERENCES integrations(id)
ON DELETE NO ACTION ON UPDATE NO ACTION;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- 1. Drop the foreign key constraint first
ALTER TABLE active_directory_settings
DROP CONSTRAINT FK_active_directory_settings_integrations;

-- 2. Then drop the column
ALTER TABLE active_directory_settings
DROP COLUMN integration_id;

-- +goose StatementEnd
