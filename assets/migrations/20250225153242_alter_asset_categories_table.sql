-- +goose Up
-- +goose StatementBegin
ALTER TABLE asset_categories DROP CONSTRAINT FK_asset_categories_asset_types;
ALTER TABLE asset_categories DROP COLUMN asset_type_id;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE asset_categories ADD asset_type_id INT NULL;
ALTER TABLE asset_categories ADD CONSTRAINT FK_asset_categories_asset_types FOREIGN KEY (asset_type_id) REFERENCES asset_types(id) ON DELETE NO ACTION ON UPDATE NO ACTION;
-- +goose StatementEnd
