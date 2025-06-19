-- +goose Up
-- +goose StatementBegin
ALTER TABLE asset_register
ADD asset_criticality_matrix_rating_scale_id INT,
    total_score FLOAT;

ALTER TABLE asset_register
ADD CONSTRAINT FK_asset_register_asset_criticality_matrix_rating_scales FOREIGN KEY (asset_criticality_matrix_rating_scale_id)
REFERENCES asset_criticality_matrix_rating_scales(id) ON DELETE NO ACTION ON UPDATE NO ACTION;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE asset_register
DROP COLUMN IF EXISTS asset_criticality_matrix_rating_scale_id,
                   total_score;
-- +goose StatementEnd
