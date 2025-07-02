-- +goose Up
-- +goose StatementBegin
CREATE TABLE compliance_report_requirements (
    id INT PRIMARY KEY IDENTITY(1,1),
    compliance_report_id INT NOT NULL,
    compliance_report_clause_id INT NOT NULL,
    title_en NVARCHAR(500) NOT NULL,
    title_ar NVARCHAR(500) NOT NULL,
    score INT NOT NULL,
    compliance_status_id INT NULL,
    created_at DATETIME NOT NULL,
    created_by_id INT NOT NULL,
    updated_at DATETIME NULL,
    updated_by_id INT NULL,
    CONSTRAINT FK_compliance_report_requirements_compliance_reports FOREIGN KEY (compliance_report_id) REFERENCES compliance_reports (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_compliance_report_requirements_compliance_report_clauses FOREIGN KEY (compliance_report_clause_id) REFERENCES compliance_report_clauses (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_compliance_report_requirements_users FOREIGN KEY (created_by_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_compliance_report_requirements_users_updated FOREIGN KEY (updated_by_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_compliance_report_requirements_compliance_statuses FOREIGN KEY (compliance_status_id) REFERENCES compliance_statuses (id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE compliance_report_requirements;
-- +goose StatementEnd
