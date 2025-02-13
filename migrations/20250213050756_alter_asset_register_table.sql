-- +goose Up
-- +goose StatementBegin
ALTER TABLE asset_register
ADD asset_sub_category_id INT NULL;
 
ALTER TABLE asset_register
ADD CONSTRAINT FK_asset_register_asset_sub_category FOREIGN KEY (asset_sub_category_id) REFERENCES asset_sub_categories(id)  ON DELETE NO ACTION ON UPDATE NO ACTION;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE asset_register
DROP CONSTRAINT FK_asset_register_asset_sub_category;

ALTER TABLE asset_register
DROP COLUMN asset_sub_category_id;
-- +goose StatementEnd
