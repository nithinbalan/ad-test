-- +goose Up
-- +goose StatementBegin
ALTER TABLE asset_register
ADD asset_criticality_matrix_id INT NULL;

ALTER TABLE asset_register
ADD CONSTRAINT FK_asset_register_asset_criticality_matrices FOREIGN KEY (asset_criticality_matrix_id) REFERENCES asset_criticality_matrices(id)
ON DELETE NO ACTION ON UPDATE NO ACTION;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE asset_register
DROP CONSTRAINT FK_asset_register_asset_criticality_matrices;

ALTER TABLE asset_register
DROP COLUMN asset_criticality_matrix_id;
-- +goose StatementEnd


