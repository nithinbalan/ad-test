-- +goose Up
-- +goose StatementBegin
CREATE TABLE audit_report_workflows (
    id INT PRIMARY KEY IDENTITY(1,1),
    audit_draft_report_id INT NULL,
    level INT NOT NULL,
    approval_types VARCHAR(255) NOT NULL DEFAULT 'anyone',
    CONSTRAINT CHK_audit_report_workflows_approval_types CHECK (approval_types IN ('anyone', 'all')),
    types VARCHAR(255) NOT NULL DEFAULT 'draft-report',
    CONSTRAINT CHK_audit_report_workflows_types CHECK (types IN ('draft-report', 'final-report')),
    audit_workflow_status_id INT NOT NULL,
    comment NVARCHAR(MAX) NULL,
    created_by_id INT NOT NULL,
    updated_by_id INT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,

    CONSTRAINT FK_audit_report_workflows_audit_draft_reports FOREIGN KEY (audit_draft_report_id) REFERENCES audit_draft_reports(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_audit_report_workflows_audit_workflow_statuses FOREIGN KEY (audit_workflow_status_id) REFERENCES audit_workflow_statuses(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_audit_report_workflows_users_1 FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_audit_report_workflows_users_2 FOREIGN KEY (updated_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.audit_report_workflows;
-- +goose StatementEnd
