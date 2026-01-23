-- +goose Up
-- +goose StatementBegin
CREATE TABLE violation_report_investigation_departments (
    id INT PRIMARY KEY IDENTITY(1,1),
    violation_report_investigation_id INT NOT NULL,
    department_id INT NOT NULL,
    created_at DATETIME NOT NULL,
    created_by_id INT NOT NULL,
    CONSTRAINT FK_violation_report_investigation_departments_violation_report_investigations FOREIGN KEY (violation_report_investigation_id) REFERENCES violation_report_investigations (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_violation_report_investigation_departments_organization_units FOREIGN KEY (department_id) REFERENCES organization_units (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_violation_report_investigation_departments_users FOREIGN KEY (created_by_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE violation_report_investigation_departments;
-- +goose StatementEnd