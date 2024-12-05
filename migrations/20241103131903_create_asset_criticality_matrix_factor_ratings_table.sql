-- +goose Up
-- +goose StatementBegin
CREATE TABLE asset_criticality_matrix_factor_ratings
(
    id INT PRIMARY KEY IDENTITY(1, 1),
    asset_criticality_matrix_id INT NOT NULL,
    title_en NVARCHAR(100) NOT NULL,
    title_ar NVARCHAR(100) NULL,
    score INT NOT NULL,
    type NVARCHAR(50) NULL,
    color_code NVARCHAR(25) NULL,
    CONSTRAINT FK_asset_criticality_matrix_factor_ratings_asset_criticality_matrices FOREIGN KEY (asset_criticality_matrix_id) REFERENCES asset_criticality_matrices (id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE asset_criticality_matrix_factor_ratings;
-- +goose StatementEnd
