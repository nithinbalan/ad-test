-- +goose Up
-- +goose StatementBegin
CREATE TABLE third_party_register
(
    id INT PRIMARY KEY IDENTITY(1, 1),
    reference_number VARCHAR(45) NOT NULL UNIQUE,
    title NVARCHAR(500) NOT NULL,
    third_party_type_id INT NOT NULL,
    country_id INT NOT NULL,
    address NVARCHAR(MAX),
    website NVARCHAR(255),
    organization_unit_id INT, 
    representative_name NVARCHAR(250),
    representative_role NVARCHAR(250),
    email NVARCHAR(150),
    phone_number VARCHAR(20),
    responsible_user_id INT NULL, 
    onboarding_date DATE,
    service_start_date DATE,
    service_end_date DATE,
    severity_level_id INT NULL,

    created_by_id INT,
    updated_by_id INT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NULL,

    CONSTRAINT FK_third_party_register_types FOREIGN KEY (third_party_type_id) REFERENCES third_party_types (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_third_party_register_countries FOREIGN KEY (country_id) REFERENCES countries (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_third_party_register_org_units FOREIGN KEY (organization_unit_id) REFERENCES organization_units (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_third_party_register_severity_levels FOREIGN KEY (severity_level_id) REFERENCES severity_levels (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_third_party_register_users_created_by FOREIGN KEY (created_by_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_third_party_register_users_updated_by FOREIGN KEY (updated_by_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_third_party_register_users_responsible FOREIGN KEY (responsible_user_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE third_party_register;
-- +goose StatementEnd
