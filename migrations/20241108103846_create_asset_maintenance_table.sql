-- +goose Up
-- +goose StatementBegin
CREATE TABLE asset_maintenance
(
    id INT PRIMARY KEY IDENTITY(1, 1),
    asset_register_id INT NOT NULL,
    type NVARCHAR(500) NOT NULL,
    performed_date DATE NOT NULL,
    title NVARCHAR(500) NOT NULL,
    description NVARCHAR(MAX) NULL,
    performed_by_id INT NOT NULL,
    outcome NVARCHAR(MAX) NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,
    created_by_id INT NOT NULL,
    updated_by_id INT NULL,
    CONSTRAINT FK_asset_maintenance_asset_register FOREIGN KEY (asset_register_id) REFERENCES asset_register (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_asset_maintenance_users_1 FOREIGN KEY (performed_by_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_asset_maintenance_users_2 FOREIGN KEY (created_by_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_asset_maintenance_users_3 FOREIGN KEY (updated_by_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS asset_maintenance;
-- +goose StatementEnd
