-- +goose Up
-- +goose StatementBegin
CREATE TABLE violation_reports
(
    id INT PRIMARY KEY IDENTITY(1, 1),
    reference_number VARCHAR(45) NOT NULL UNIQUE,
    reported_at DATE NOT NULL,
    violation_report_channel_id INT NOT NULL,
    id_of_reporter INT NULL,
    name_of_reporter NVARCHAR(MAX) NULL,
    phone_of_reporter NVARCHAR(MAX) NULL,
    title_en NVARCHAR(MAX) NOT NULL,
    title_ar NVARCHAR(MAX) NOT NULL,
    reported_entity_id INT NOT NULL,
    description_en NVARCHAR(MAX) NULL,
    description_ar NVARCHAR(MAX) NULL,
    violation_report_category_id INT NOT NULL,
    violation_report_sub_category_id INT NULL,
    violation_report_reason_id INT NOT NULL,
    severity_level_id INT NULL,
    evidence_description NVARCHAR(MAX) NULL,
    violation_report_status_id INT NOT NULL,
    created_by_id INT NOT NULL,
    updated_by_id INT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,
    CONSTRAINT FK_violation_reports_violation_report_channels FOREIGN KEY (violation_report_channel_id) REFERENCES violation_report_channels (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_violation_reports_violation_report_categories FOREIGN KEY (violation_report_category_id) REFERENCES violation_report_categories (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_violation_reports_violation_report_sub_categories FOREIGN KEY (violation_report_sub_category_id) REFERENCES violation_report_sub_categories (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_violation_reports_violation_report_reasons FOREIGN KEY (violation_report_reason_id) REFERENCES violation_report_reasons (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_violation_reports_violation_report_statuses FOREIGN KEY (violation_report_status_id) REFERENCES violation_report_statuses (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_violation_reports_severity_levels FOREIGN KEY (severity_level_id) REFERENCES severity_levels (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_violation_reports_organization_units FOREIGN KEY (reported_entity_id) REFERENCES organization_units (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_violation_reports_users_1 FOREIGN KEY (created_by_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_violation_reports_users_2 FOREIGN KEY (updated_by_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE violation_reports;

-- +goose StatementEnd