-- +goose Up
-- +goose StatementBegin
CREATE TABLE grc_report_compliance_departments (
    id INT PRIMARY KEY IDENTITY(1,1),
    grc_report_id INT NOT NULL,
    grc_report_compliance_summary_id INT NOT NULL,
    organization_unit_id INT NOT NULL,
    compliance_score DECIMAL(5,2) NULL,
    compliance_status_id INT NOT NULL,
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL,
    updated_by_id INT NULL,
    updated_at DATETIME NULL,

    CONSTRAINT FK_grc_report_compliance_departments_grc_report FOREIGN KEY (grc_report_id) REFERENCES grc_reports(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_grc_rcd_grc_report_compliance_summary FOREIGN KEY (grc_report_compliance_summary_id) REFERENCES grc_report_compliance_summary(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_grc_report_compliance_departments_organization_unit FOREIGN KEY (organization_unit_id) REFERENCES organization_units(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_grc_report_compliance_departments_user_1 FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_grc_report_compliance_departments_user_2 FOREIGN KEY (updated_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_grc_report_compliance_departments_compliance_status FOREIGN KEY (compliance_status_id) REFERENCES compliance_statuses(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE grc_report_compliance_departments;
-- +goose StatementEnd
