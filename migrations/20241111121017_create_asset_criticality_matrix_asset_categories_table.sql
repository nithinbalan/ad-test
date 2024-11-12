-- +goose Up
-- +goose StatementBegin
CREATE TABLE asset_criticality_matrix_asset_categories
(
    id INT PRIMARY KEY IDENTITY(1, 1),
    asset_criticality_matrix_id INT NOT NULL,
    asset_category_id INT NOT NULL,
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL,
    CONSTRAINT FK_asset_criticality_matrix_asset_categories_matrix FOREIGN KEY (asset_criticality_matrix_id) REFERENCES asset_criticality_matrices(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_asset_criticality_matrix_asset_categories_category FOREIGN KEY (asset_category_id) REFERENCES asset_categories(id) ON DELETE NO ACTION ON UPDATE NO ACTION
    CONSTRAINT FK_asset_criticality_matrix_users_1 FOREIGN KEY (created_by_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION,

);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE asset_criticality_matrix_asset_categories;

-- +goose StatementEnd
