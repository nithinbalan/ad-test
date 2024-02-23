-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
    id INT PRIMARY KEY IDENTITY(1,1),
    organization_unit_1_id INT,
    organization_unit_2_id INT,
    organization_unit_3_id INT,
    organization_unit_4_id INT,
    organization_unit_5_id INT,
    first_name_en NVARCHAR(100) NOT NULL,
    last_name_en NVARCHAR(100) NOT NULL,
    first_name_ar NVARCHAR(100) NOT NULL,
    last_name_ar NVARCHAR(100) NOT NULL,
    email VARCHAR(100) NOT NULL,
    password_hash VARCHAR(MAX) NOT NULL,
    profile_image_path VARCHAR(255),
    created_at DATETIME NOT NULL DEFAULT GETDATE(),
    updated_at DATETIME NOT NULL DEFAULT GETDATE(),
    CONSTRAINT FK_users_organization_unit_1 FOREIGN KEY (organization_unit_1_id) REFERENCES organization_units(id),
    CONSTRAINT FK_users_organization_unit_2 FOREIGN KEY (organization_unit_2_id) REFERENCES organization_units(id),
    CONSTRAINT FK_users_organization_unit_3 FOREIGN KEY (organization_unit_3_id) REFERENCES organization_units(id),
    CONSTRAINT FK_users_organization_unit_4 FOREIGN KEY (organization_unit_4_id) REFERENCES organization_units(id),
    CONSTRAINT FK_users_organization_unit_5 FOREIGN KEY (organization_unit_5_id) REFERENCES organization_units(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.users;
-- +goose StatementEnd
