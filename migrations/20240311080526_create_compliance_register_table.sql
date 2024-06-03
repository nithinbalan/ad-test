-- +goose Up
-- +goose StatementBegin
CREATE TABLE compliance_register (
    id INT PRIMARY KEY IDENTITY(1, 1),
    compliance_category_id INT NOT NULL,
    compliance_status_id INT NOT NULL,
    responsible_unit_id INT NOT NULL,
    responsible_user_id INT NOT NULL,
    review_frequency_id INT NULL,
    review_responsible_user_id INT NULL,
    severity_level_id INT NULL,
    reference_number VARCHAR(45) NULL,
    standard_code VARCHAR(45) NULL,
    title_en NVARCHAR(255) NOT NULL,
    title_ar NVARCHAR(255) NOT NULL,
    description_en NVARCHAR(MAX) NULL,
    description_ar NVARCHAR(MAX) NULL,
    version VARCHAR(45) NULL,
    deadline DATE NULL,
    issue_date DATE NULL,
    expiry_date DATE NULL,
    next_review_date DATE NULL,
    last_reviewed_date DATE NULL,
    last_reviewed_user_id INT NULL,
    last_review_note NVARCHAR(MAX) NULL,
    compliance_review_status_id INT NULL,
    created_at DATETIME NOT NULL ,
    updated_at DATETIME NULL,
    created_by INT NOT NULL,
    updated_by INT NULL,
    lifecycle_status_id INT NOT NULL,
    CONSTRAINT FK_compliance_register_compliance_categories FOREIGN KEY (compliance_category_id) REFERENCES compliance_categories (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_compliance_register_severity_levels FOREIGN KEY (severity_level_id) REFERENCES severity_levels (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_compliance_register_compliance_statuses FOREIGN KEY (compliance_status_id) REFERENCES compliance_statuses (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_compliance_register_organization_units FOREIGN KEY (responsible_unit_id) REFERENCES organization_units (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_compliance_register_users_1 FOREIGN KEY (responsible_user_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_compliance_register_review_frequencies FOREIGN KEY (review_frequency_id) REFERENCES review_frequencies (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_compliance_register_users_2 FOREIGN KEY (review_responsible_user_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_compliance_register_users_3 FOREIGN KEY (created_by) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_compliance_register_users_4 FOREIGN KEY (updated_by) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_compliance_register_lifecycle_statuses FOREIGN KEY (lifecycle_status_id) REFERENCES lifecycle_statuses (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_compliance_register_users_5 FOREIGN KEY (last_reviewed_user_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_compliance_register_compliance_review_statuses FOREIGN KEY (compliance_review_status_id) REFERENCES compliance_review_statuses (id) ON DELETE NO ACTION ON UPDATE NO ACTION
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE compliance_register;

-- +goose StatementEnd
