-- +goose Up
-- +goose StatementBegin
CREATE TABLE compliance_requirements (
    id INT PRIMARY KEY IDENTITY(1, 1),
    compliance_clause_id INT NOT NULL,
    title VARCHAR(MAX) NOT NULL,
    compliance_requirement_category_id INT NOT NULL,
    severity_level_id INT NOT NULL,
    responsible_user_id INT NOT NULL,
    deadline DATE NULL,
    compliance_status_id INT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT GETDATE(),
    updated_at DATETIME NOT NULL DEFAULT GETDATE(),
    created_by INT NOT NULL,
    updated_by INT NULL,
    CONSTRAINT FK_compliance_requirements_compliance_clauses FOREIGN KEY (compliance_clause_id) REFERENCES compliance_clauses(id),
    CONSTRAINT FK_compliance_requirements_compliance_requirement_categories FOREIGN KEY (compliance_requirement_category_id) REFERENCES compliance_requirement_categories(id),
    CONSTRAINT FK_compliance_requirements_severity_levels FOREIGN KEY (severity_level_id) REFERENCES severity_levels(id),
    CONSTRAINT FK_compliance_requirements_users1 FOREIGN KEY (responsible_user_id) REFERENCES users(id),
    CONSTRAINT FK_compliance_requirements_compliance_statuses FOREIGN KEY (compliance_status_id) REFERENCES compliance_statuses(id),
    CONSTRAINT FK_compliance_requirements_users2 FOREIGN KEY (created_by) REFERENCES users(id),
    CONSTRAINT FK_compliance_requirements_users3 FOREIGN KEY (updated_by) REFERENCES users(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE compliance_requirements;
-- +goose StatementEnd
