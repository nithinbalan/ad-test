-- +goose Up
-- +goose StatementBegin
ALTER TABLE asset_register
ADD tenable_uuid VARCHAR(200) NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE asset_register
DROP COLUMN tenable_uuid;
-- +goose StatementEnd
