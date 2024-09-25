-- +goose Up
-- +goose StatementBegin
CREATE TABLE audit_report_workflows
(
    id INT PRIMARY KEY IDENTITY(1,1),
    audit_draft_report_id INT NULL,
    level INT NOT NULL,
    approval_type VARCHAR(255) NOT NULL DEFAULT 'anyone',
    CONSTRAINT CHK_audit_report_workflows_approval_type CHECK (approval_type IN ('anyone', 'all')),
    type VARCHAR(50) NOT NULL DEFAULT 'draft',
    CONSTRAINT CHK_audit_report_workflows_type CHECK (type IN ('draft', 'final')),
    workflow_status_id INT NOT NULL,
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL,

    CONSTRAINT FK_audit_report_workflows_audit_draft_reports FOREIGN KEY (audit_draft_report_id) REFERENCES audit_draft_reports(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_audit_report_workflows_workflow_statuses FOREIGN KEY (workflow_status_id) REFERENCES workflow_statuses(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_audit_report_workflows_users_1 FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.audit_report_workflows;
-- +goose StatementEnd
