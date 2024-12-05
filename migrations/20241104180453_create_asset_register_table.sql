-- +goose Up
-- +goose StatementBegin
CREATE TABLE asset_register (
    id INT PRIMARY KEY IDENTITY(1,1),
    reference_number VARCHAR(45) NOT NULL UNIQUE,
    title_en NVARCHAR(500) NOT NULL,
    title_ar NVARCHAR(500) NOT NULL,
    asset_type_id INT NOT NULL,
    asset_category_id INT NOT NULL,
    asset_status_id INT NOT NULL,
    responsible_user_id INT NOT NULL,
    responsible_unit_id INT NOT NULL,
    location NVARCHAR(255),
    description NVARCHAR(MAX),
    acquisition_date DATE,
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL,
    updated_by_id INT,
    updated_at DATETIME,

    CONSTRAINT FK_asset_register_asset_types FOREIGN KEY (asset_type_id) REFERENCES asset_types(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_asset_register_asset_categories FOREIGN KEY (asset_category_id) REFERENCES asset_categories(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_asset_register_asset_statuses FOREIGN KEY (asset_status_id) REFERENCES asset_statuses(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_asset_register_organization_units FOREIGN KEY (responsible_unit_id) REFERENCES organization_units(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_asset_register_users_1 FOREIGN KEY (responsible_user_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_asset_register_users_2 FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_asset_register_users_3 FOREIGN KEY (updated_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS asset_register;
-- +goose StatementEnd
