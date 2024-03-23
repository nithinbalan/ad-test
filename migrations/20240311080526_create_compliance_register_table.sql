-- +goose Up
-- +goose StatementBegin
CREATE TABLE compliance_register (
    id INT PRIMARY KEY IDENTITY(1,1),
    compliance_category_id INT NOT NULL,
    compliance_status_id INT NOT NULL,
    responsible_unit_id INT NOT NULL,
    responsible_user_id INT NOT NULL,
    review_frequency_id INT NOT NULL,
    review_responsible_user_id INT NULL,
    severity_level_id INT NULL,
    reference_number VARCHAR(45) NULL,
    standard_code VARCHAR(45) NULL,
    title_en VARCHAR(255) NOT NULL,
    title_ar NVARCHAR(255) NOT NULL,
    description_en VARCHAR(MAX) NULL,
    description_ar NVARCHAR(MAX) NULL,
    version VARCHAR(45) NULL,
    deadline DATE NULL,
    issue_date DATE NULL,
    expiry_date DATE NULL,
    next_review_date DATE NULL,
    created_at DATETIME NOT NULL DEFAULT GETDATE(),
    updated_at DATETIME NOT NULL DEFAULT GETDATE(),
    created_by INT NOT NULL,
    updated_by INT NULL,
    lifecycle_status_id INT NOT NULL,
    CONSTRAINT fk_compliance_registers_compliance_categories1 FOREIGN KEY (compliance_category_id) REFERENCES compliance_categories (id),
    CONSTRAINT fk_compliance_registers_severity_levels1 FOREIGN KEY (severity_level_id) REFERENCES severity_levels (id),
    CONSTRAINT fk_compliance_registers_compliance_statuses1 FOREIGN KEY (compliance_status_id) REFERENCES compliance_statuses (id),
    CONSTRAINT fk_compliance_registers_organization_units1 FOREIGN KEY (responsible_unit_id) REFERENCES organization_units (id),
    CONSTRAINT fk_compliance_registers_users1 FOREIGN KEY (responsible_user_id) REFERENCES users (id),
    CONSTRAINT fk_compliance_registers_review_frequencies1 FOREIGN KEY (review_frequency_id) REFERENCES review_frequencies (id),
    CONSTRAINT fk_compliance_registers_users2 FOREIGN KEY (review_responsible_user_id) REFERENCES users (id),
    CONSTRAINT fk_compliance_registers_users3 FOREIGN KEY (created_by) REFERENCES users (id),
    CONSTRAINT fk_compliance_registers_users4 FOREIGN KEY (updated_by) REFERENCES users (id),
    CONSTRAINT fk_compliance_registers_lifecycle_statuses1 FOREIGN KEY (lifecycle_status_id) REFERENCES lifecycle_statuses (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE compliance_register;
-- +goose StatementEnd
