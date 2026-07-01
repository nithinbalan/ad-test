-- +goose Up
-- +goose StatementBegin
ALTER TABLE asset_register
ADD
    technical_owner_id INT NULL,
    administrative_owner_id INT NULL,
    hostname NVARCHAR(255) NULL,
    ip_address NVARCHAR(100) NULL,
    mac_address NVARCHAR(100) NULL,
    fqdn NVARCHAR(255) NULL,
    device_category_id INT NULL,
    criticality_classification_id INT NULL,
    confidentiality_classification_id INT NULL,
    region NVARCHAR(255) NULL,
    assigned_user_id INT NULL,
    operating_system NVARCHAR(255) NULL,
    asset_environment_id INT NULL,
    inventory_time DATETIME NULL,
    discovery_time DATETIME NULL,
    last_response_time DATETIME NULL;

ALTER TABLE asset_register ADD CONSTRAINT FK_asset_register_user_4 FOREIGN KEY (technical_owner_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION;
ALTER TABLE asset_register ADD CONSTRAINT FK_asset_register_user_5 FOREIGN KEY (administrative_owner_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION;
ALTER TABLE asset_register ADD CONSTRAINT FK_asset_register_device_categories FOREIGN KEY (device_category_id) REFERENCES asset_device_categories(id) ON DELETE NO ACTION ON UPDATE NO ACTION;
ALTER TABLE asset_register ADD CONSTRAINT FK_asset_register_criticality_classifications FOREIGN KEY (criticality_classification_id) REFERENCES asset_criticality_matrix_factor_ratings(id) ON DELETE NO ACTION ON UPDATE NO ACTION;
ALTER TABLE asset_register ADD CONSTRAINT FK_asset_register_confidentiality_classifications FOREIGN KEY (confidentiality_classification_id) REFERENCES confidentiality_classifications(id) ON DELETE NO ACTION ON UPDATE NO ACTION;
ALTER TABLE asset_register ADD CONSTRAINT FK_asset_register_asset_environments FOREIGN KEY (asset_environment_id) REFERENCES asset_environments(id) ON DELETE NO ACTION ON UPDATE NO ACTION;
ALTER TABLE asset_register ADD CONSTRAINT FK_asset_register_assigned_user_6 FOREIGN KEY (assigned_user_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE asset_register DROP CONSTRAINT FK_asset_register_user_4;
ALTER TABLE asset_register DROP CONSTRAINT FK_asset_register_user_5;
ALTER TABLE asset_register DROP CONSTRAINT FK_asset_register_device_categories;
ALTER TABLE asset_register DROP CONSTRAINT FK_asset_register_criticality_classifications;
ALTER TABLE asset_register DROP CONSTRAINT FK_asset_register_confidentiality_classifications;
ALTER TABLE asset_register DROP CONSTRAINT FK_asset_register_asset_environments;
ALTER TABLE asset_register DROP CONSTRAINT FK_asset_register_assigned_user_6;
ALTER TABLE asset_register
DROP COLUMN
    technical_owner_id,
    administrative_owner_id,
    hostname,
    ip_address,
    mac_address,
    fqdn,
    device_category_id,
    criticality_classification_id,
    confidentiality_classification_id,
    region,
    assigned_user_id,
    operating_system,
    asset_environment_id,
    inventory_time,
    discovery_time,
    last_response_time;
-- +goose StatementEnd