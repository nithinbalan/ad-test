-- +goose Up
-- +goose StatementBegin

ALTER TABLE third_party_register
ADD responsible_unit_id INT NULL;

ALTER TABLE third_party_register
ADD CONSTRAINT FK_third_party_register_organization_units FOREIGN KEY (responsible_unit_id) REFERENCES organization_units(id) ON DELETE NO ACTION ON UPDATE NO ACTION;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE third_party_register
DROP CONSTRAINT FK_third_party_register_organization_units;

ALTER TABLE third_party_register
DROP COLUMN responsible_unit_id;
-- +goose StatementEnd
