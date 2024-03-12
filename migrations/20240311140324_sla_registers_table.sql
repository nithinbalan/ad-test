-- +goose Up
-- +goose StatementBegin
CREATE TABLE sla_registers (
    id INT PRIMARY KEY IDENTITY(1,1),
    sla_category_id INT,
    severity_level_id INT,
    responsible_unit_id INT,
    responsible_user_id INT,
    lifecycle_status_id INT,
    title_en VARCHAR(100) NOT NULL,
    title_ar NVARCHAR(100) NOT NULL,
    description VARCHAR(MAX) NULL,
    client VARCHAR(100) NULL,
    created_by INT,
    updated_by INT NULL,
    created_at DATETIME NOT NULL DEFAULT GETDATE(),
    updated_at DATETIME NOT NULL DEFAULT GETDATE(),
    CONSTRAINT FK_sla_registers_sla_categories_key1 FOREIGN KEY (sla_category_id) REFERENCES sla_categories(id),
    CONSTRAINT FK_sla_registers_severity_levels_key1 FOREIGN KEY (severity_level_id) REFERENCES severity_levels(id),
    CONSTRAINT FK_sla_registers_responsible_units_key1 FOREIGN KEY (responsible_unit_id) REFERENCES organization_units(id),
    CONSTRAINT FK_sla_registers_responsible_users_key1 FOREIGN KEY (responsible_user_id) REFERENCES users(id),
    CONSTRAINT FK_sla_registers_lifecycle_statuses_key1 FOREIGN KEY (lifecycle_status_id) REFERENCES lifecycle_statuses(id)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.sla_registers;
-- +goose StatementEnd
