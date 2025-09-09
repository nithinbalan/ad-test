-- +goose Up
-- +goose StatementBegin
ALTER TABLE roles
ADD license_type_id INT NULL;

ALTER TABLE roles
ADD CONSTRAINT FK_roles_license_types
FOREIGN KEY (license_type_id) REFERENCES license_types(id) ON DELETE NO ACTION ON UPDATE NO ACTION;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Remove the foreign key constraint for license_type_id
ALTER TABLE roles
DROP CONSTRAINT FK_roles_license_types;

-- Drop the license_type_id
ALTER TABLE roles
DROP COLUMN license_type_id;
-- +goose StatementEnd
