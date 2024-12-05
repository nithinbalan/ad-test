-- +goose Up
-- +goose StatementBegin
CREATE TABLE asset_criticality_matrix_factors
(
    id INT PRIMARY KEY IDENTITY(1, 1),
    asset_criticality_matrix_id INT NOT NULL,
    title_en NVARCHAR(100) NOT NULL,
    title_ar NVARCHAR(100) NULL,
    description_en NVARCHAR(MAX) NULL,
    description_ar NVARCHAR(MAX) NULL,
    weightage FLOAT NOT NULL,    
    CONSTRAINT FK_asset_criticality_matrix_factors_asset_criticality_matrices FOREIGN KEY (asset_criticality_matrix_id) REFERENCES asset_criticality_matrices (id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE asset_criticality_matrix_factors;
-- +goose StatementEnd
