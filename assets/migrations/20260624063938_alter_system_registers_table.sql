-- +goose Up
-- +goose StatementBegin
ALTER TABLE system_registers DROP CONSTRAINT FK_system_registers_impact_categories;
ALTER TABLE system_registers DROP COLUMN impact_type_id;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE system_registers ADD impact_type_id INT NULL;
ALTER TABLE system_registers ADD CONSTRAINT FK_system_registers_impact_categories FOREIGN KEY (impact_type_id) REFERENCES impact_categories(id) ON DELETE NO ACTION ON UPDATE NO ACTION;
-- +goose StatementEnd
