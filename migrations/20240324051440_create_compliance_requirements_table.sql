-- +goose Up
-- +goose StatementBegin
CREATE TABLE compliance_requirements (
    id INT PRIMARY KEY IDENTITY(1, 1),
    compliance_clause_id INT NOT NULL,
    title NVARCHAR(MAX) NOT NULL,
    compliance_requirement_category_id INT NOT NULL,
    severity_level_id INT NOT NULL,
    responsible_user_id INT NOT NULL,
    deadline DATE NULL,
    compliance_status_id INT NOT NULL,
    created_at DATETIME NOT NULL ,
    updated_at DATETIME NULL,
    created_by_id INT NOT NULL,
    updated_by_id INT NULL,
    CONSTRAINT FK_compliance_requirements_compliance_clauses FOREIGN KEY (compliance_clause_id) REFERENCES compliance_clauses(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_compliance_requirements_compliance_requirement_categories FOREIGN KEY (compliance_requirement_category_id) REFERENCES compliance_requirement_categories(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_compliance_requirements_severity_levels FOREIGN KEY (severity_level_id) REFERENCES severity_levels(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_compliance_requirements_users_1 FOREIGN KEY (responsible_user_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_compliance_requirements_compliance_statuses FOREIGN KEY (compliance_status_id) REFERENCES compliance_statuses(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_compliance_requirements_users_2 FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_compliance_requirements_users_3 FOREIGN KEY (updated_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE compliance_requirements;
-- +goose StatementEnd
