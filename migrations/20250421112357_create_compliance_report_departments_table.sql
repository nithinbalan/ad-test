-- +goose Up
-- +goose StatementBegin
CREATE TABLE compliance_report_departments (
    id INT PRIMARY KEY IDENTITY(1,1),
    compliance_report_id INT NOT NULL,
    organization_unit_id INT NOT NULL,
    compliance_percentage DECIMAL(10,2) NOT NULL DEFAULT 0.00,
    created_at DATETIME NOT NULL,
    created_by_id INT NOT NULL,
    CONSTRAINT FK_compliance_report_departments_compliance_reports FOREIGN KEY (compliance_report_id) REFERENCES compliance_reports (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_compliance_report_departments_organization_units FOREIGN KEY (organization_unit_id) REFERENCES organization_units (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_compliance_report_departments_users FOREIGN KEY (created_by_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE compliance_report_departments;
-- +goose StatementEnd