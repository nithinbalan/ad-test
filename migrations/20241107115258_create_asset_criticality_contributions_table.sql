-- +goose Up
-- +goose StatementBegin
CREATE TABLE asset_criticality_contributions (
    id INT PRIMARY KEY IDENTITY(1,1),
    asset_register_id INT NOT NULL,
    asset_criticality_matrix_factor_id INT NOT NULL,
    asset_criticality_matrix_factor_rating_id INT NOT NULL,
    score FLOAT NOT NULL,
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL,
    updated_by_id INT,
    updated_at DATETIME,

    CONSTRAINT FK_asset_criticality_contributions_asset_register FOREIGN KEY (asset_register_id) REFERENCES asset_register(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_asset_criticality_contributions_asset_criticality_matrix_factor FOREIGN KEY (asset_criticality_matrix_factor_id) REFERENCES asset_criticality_matrix_factors(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_asset_criticality_contributions_asset_criticality_matrix_factor_rating FOREIGN KEY (asset_criticality_matrix_factor_rating_id) REFERENCES asset_criticality_matrix_factor_ratings(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_asset_criticality_contributions_users_1 FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_asset_criticality_contributions_users_2 FOREIGN KEY (updated_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS asset_criticality_contributions;
-- +goose StatementEnd
