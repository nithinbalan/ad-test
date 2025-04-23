-- +goose Up
-- +goose StatementBegin
CREATE TABLE compliance_report_clauses (
    id INT PRIMARY KEY IDENTITY(1,1),
    compliance_report_id INT NOT NULL,
    compliance_report_department_id INT NOT NULL,
    clause_no NVARCHAR(45) NOT NULL,
    clause_title_en NVARCHAR(500) NOT NULL,
    clause_title_ar NVARCHAR(500) NULL,
    description_en NVARCHAR(MAX) NULL,
    description_ar NVARCHAR(MAX) NULL,
    created_at DATETIME NOT NULL,
    created_by_id INT NOT NULL,
    updated_at DATETIME NULL,
    updated_by_id INT NULL,
    CONSTRAINT FK_compliance_report_clauses_compliance_reports FOREIGN KEY (compliance_report_id) REFERENCES compliance_reports (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_compliance_report_clauses_users FOREIGN KEY (created_by_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_compliance_report_clauses_users_updated FOREIGN KEY (updated_by_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_compliance_report_clauses_compliance_report_departments FOREIGN KEY (compliance_report_department_id) REFERENCES compliance_report_departments (id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE compliance_report_clauses;
-- +goose StatementEnd
