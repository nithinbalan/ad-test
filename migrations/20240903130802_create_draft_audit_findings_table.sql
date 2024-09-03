-- +goose Up
-- +goose StatementBegin
CREATE TABLE draft_audit_findings (
    id INT PRIMARY KEY IDENTITY(1, 1),
    audit_id INT NOT NULL,
    title NVARCHAR(255) NOT NULL,
    description NVARCHAR(MAX) NULL,
    audit_checklist_id INT NULL,
    auditable_item_id INT NULL,
    audit_finding_category_id INT NOT NULL,
    severity_level_id INT NULL,
    recommendations NVARCHAR(MAX),
    draft_audit_finding_status_id INT NOT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,
    created_by_id INT NOT NULL,
    updated_by_id INT NULL,

    CONSTRAINT FK_draft_audit_findings_audits FOREIGN KEY (audit_id) REFERENCES audits(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_draft_audit_findings_audit_checklists FOREIGN KEY (audit_checklist_id) REFERENCES audit_checklists(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_draft_audit_findings_auditable_items FOREIGN KEY (auditable_item_id) REFERENCES auditable_items(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_draft_audit_findings_categories FOREIGN KEY (audit_finding_category_id) REFERENCES audit_finding_categories(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_draft_audit_findings_severity_levels FOREIGN KEY (severity_level_id) REFERENCES severity_levels(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_draft_audit_findings_draft_audit_finding_statuses FOREIGN KEY (draft_audit_finding_status_id) REFERENCES draft_audit_finding_statuses(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_draft_audit_findings_users_1 FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_draft_audit_findings_users_2 FOREIGN KEY (updated_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE draft_audit_findings;
-- +goose StatementEnd
