-- +goose Up
-- +goose StatementBegin

ALTER TABLE organization_units
ADD organization_unit_head_id INT NULL;

ALTER TABLE organization_units
ADD CONSTRAINT FK_organization_units_head_user
FOREIGN KEY (organization_unit_head_id)
REFERENCES users(id);

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

ALTER TABLE organization_units
DROP CONSTRAINT FK_organization_units_head_user;

ALTER TABLE organization_units
DROP COLUMN organization_unit_head_id;

-- +goose StatementEnd