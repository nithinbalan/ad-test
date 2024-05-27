-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
    id INT PRIMARY KEY IDENTITY(1,1),
    organization_unit_1_id INT,
    organization_unit_2_id INT,
    organization_unit_3_id INT,
    organization_unit_4_id INT,
    organization_unit_5_id INT,
    user_status_id INT,
    designation_id INT,
    first_name_en NVARCHAR(255) NOT NULL,
    last_name_en NVARCHAR(255) NOT NULL,
    first_name_ar NVARCHAR(255) NOT NULL,
    last_name_ar NVARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    mobile_number VARCHAR(45),
    profile_image_path VARCHAR(500),
    last_login_at DATETIME NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,
    created_by INT,
    updated_by INT NULL,
    CONSTRAINT FK_users_organization_units_1 FOREIGN KEY (organization_unit_1_id) REFERENCES organization_units(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_users_organization_units_2 FOREIGN KEY (organization_unit_2_id) REFERENCES organization_units(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_users_organization_units_3 FOREIGN KEY (organization_unit_3_id) REFERENCES organization_units(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_users_organization_units_4 FOREIGN KEY (organization_unit_4_id) REFERENCES organization_units(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_users_organization_units_5 FOREIGN KEY (organization_unit_5_id) REFERENCES organization_units(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_users_user_statuses FOREIGN KEY (user_status_id) REFERENCES user_statuses(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_users_designations FOREIGN KEY (designation_id) REFERENCES designations(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_users_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_users_updated_by FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.users;
-- +goose StatementEnd
