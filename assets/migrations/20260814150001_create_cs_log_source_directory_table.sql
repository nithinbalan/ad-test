-- +goose Up
-- +goose StatementBegin
CREATE TABLE cs_log_source_directory (
    id INT PRIMARY KEY IDENTITY(1, 1),
    reference_number NVARCHAR(45) NOT NULL,
    entity_name NVARCHAR(255) NOT NULL,
    data_entity_type_id INT NULL,
    location NVARCHAR(500) NULL,
    responsible_entity_id INT NULL,
    beneficiary_entity_id INT NULL,
    contact_name NVARCHAR(255) NULL,
    contact_mobile VARCHAR(50) NULL,
    contact_email VARCHAR(255) NULL,
    website NVARCHAR(500) NULL,
    notes NVARCHAR(MAX) NULL,
    created_by_id INT NOT NULL,
    updated_by_id INT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,

    CONSTRAINT UQ_cs_log_source_directory_reference_number UNIQUE (reference_number),

    CONSTRAINT FK_cs_log_source_directory_cs_log_data_entity_types
        FOREIGN KEY (data_entity_type_id) REFERENCES cs_log_data_entity_types(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_cs_log_source_directory_organization_units_1
        FOREIGN KEY (responsible_entity_id) REFERENCES organization_units(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_cs_log_source_directory_organization_units_2
        FOREIGN KEY (beneficiary_entity_id) REFERENCES organization_units(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_cs_log_source_directory_users_1
        FOREIGN KEY (created_by_id) REFERENCES users(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_cs_log_source_directory_users_2
        FOREIGN KEY (updated_by_id) REFERENCES users(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE cs_log_source_directory;
-- +goose StatementEnd
