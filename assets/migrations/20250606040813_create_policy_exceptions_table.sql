-- +goose Up
-- +goose StatementBegin
CREATE TABLE policy_exceptions (
    id INT PRIMARY KEY IDENTITY(1,1),
    document_id INT NOT NULL,
    document_version_id INT NOT NULL,
    document_version_clause_id INT NULL,
    document_requirement_group_id INT NULL,
    reference_number VARCHAR(45) NOT NULL,
    title_en NVARCHAR(500) NOT NULL,
    title_ar NVARCHAR(500) NOT NULL,
    requirements NVARCHAR(MAX) NULL,
    justification NVARCHAR(MAX) NULL,
    start_date DATE NOT NULL,
    end_date DATE NULL,
    requested_by_id INT NOT NULL,
    organization_unit_id INT NOT NULL,
    workflow_status_id INT NULL,
    created_at DATETIME NOT NULL,
    created_by_id INT NOT NULL,
    updated_at DATETIME NULL,
    updated_by_id INT NULL,
 
    CONSTRAINT FK_policy_exceptions_policies FOREIGN KEY (document_id) REFERENCES documents (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_policy_exceptions_document_versions FOREIGN KEY (document_version_id) REFERENCES document_versions (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_policy_exceptions_document_version_clauses FOREIGN KEY (document_version_clause_id) REFERENCES document_version_clauses (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_policy_exceptions_document_requirement_groups FOREIGN KEY (document_requirement_group_id) REFERENCES document_requirement_groups (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_policy_exceptions_workflow_statuses FOREIGN KEY (workflow_status_id) REFERENCES workflow_statuses (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_policy_exceptions_users_requested FOREIGN KEY (requested_by_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_policy_exceptions_org_units FOREIGN KEY (organization_unit_id) REFERENCES organization_units (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_policy_exceptions_users_created FOREIGN KEY (created_by_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_policy_exceptions_users_updated FOREIGN KEY (updated_by_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd
 
-- +goose Down
-- +goose StatementBegin
DROP TABLE policy_exceptions;
-- +goose StatementEnd