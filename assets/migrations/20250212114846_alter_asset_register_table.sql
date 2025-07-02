-- +goose Up
-- +goose StatementBegin
ALTER TABLE asset_register
ADD assignment_date DATE NULL;

ALTER TABLE asset_register
ADD warranty_expiry_date DATE NULL;

ALTER TABLE asset_register
ADD manufacturer VARCHAR(255) NULL;

ALTER TABLE asset_register
ADD model VARCHAR(255) NULL;

ALTER TABLE asset_register
ADD serial_number VARCHAR(255) NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE asset_register
DROP COLUMN assignment_date;

ALTER TABLE asset_register
DROP COLUMN warranty_expiry_date;

ALTER TABLE asset_register
DROP COLUMN manufacturer;

ALTER TABLE asset_register
DROP COLUMN model;

ALTER TABLE asset_register
DROP COLUMN serial_number;
-- +goose StatementEnd
