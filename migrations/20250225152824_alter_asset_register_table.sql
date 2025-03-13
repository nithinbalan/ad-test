-- +goose Up
-- +goose StatementBegin
ALTER TABLE asset_register DROP CONSTRAINT FK_asset_register_asset_types;
ALTER TABLE asset_register DROP COLUMN asset_type_id;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE asset_register ADD asset_type_id INT NOT NULL;
ALTER TABLE asset_register ADD CONSTRAINT FK_asset_register_asset_types FOREIGN KEY (asset_type_id) REFERENCES asset_types(id) ON DELETE NO ACTION ON UPDATE NO ACTION;
-- +goose StatementEnd
