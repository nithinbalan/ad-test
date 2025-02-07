-- +goose Up
-- +goose StatementBegin
ALTER TABLE integrations
DROP CONSTRAINT FK_integrations_integration_categories;

ALTER TABLE integrations
ADD CONSTRAINT FK_integrations_integration_categories
FOREIGN KEY (integration_category_id) REFERENCES integration_categories(id) ON DELETE NO ACTION ON UPDATE NO ACTION;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE integrations
DROP CONSTRAINT FK_integrations_integration_categories;

ALTER TABLE integrations
ADD CONSTRAINT FK_integrations_integration_categories
FOREIGN KEY (integration_category_id) REFERENCES integrations(id) ON DELETE NO ACTION ON UPDATE NO ACTION;
-- +goose StatementEnd
