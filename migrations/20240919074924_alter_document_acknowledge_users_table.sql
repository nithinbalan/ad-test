-- +goose Up
-- +goose StatementBegin
ALTER TABLE document_acknowledge_users
ADD organization_unit_id INT NULL;

ALTER TABLE document_acknowledge_users
ADD CONSTRAINT FK_document_acknowledge_users_organization_units FOREIGN KEY (organization_unit_id) 
REFERENCES organization_units(id) ON DELETE NO ACTION ON UPDATE NO ACTION;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE document_acknowledge_users
DROP CONSTRAINT FK_document_acknowledge_users_organization_units;

ALTER TABLE document_acknowledge_users
DROP COLUMN organization_unit_id;
-- +goose StatementEnd
