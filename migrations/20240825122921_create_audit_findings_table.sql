-- +goose Up
-- +goose StatementBegin
CREATE TABLE audit_findings (
    id INT PRIMARY KEY IDENTITY(1, 1),
    audit_id INT NOT NULL,
    auditable_item_id INT NOT NULL,
    audit_finding_category_id INT NOT NULL,
    severity_level_id INT NOT NULL,
    recommendations NVARCHAR(MAX),
    audit_finding_status_id INT NOT NULL,

    CONSTRAINT FK_audit_findings_audits FOREIGN KEY (audit_id) REFERENCES audits(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_audit_findings_auditable_items FOREIGN KEY (auditable_item_id) REFERENCES auditable_items(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_audit_findings_categories FOREIGN KEY (audit_finding_category_id) REFERENCES audit_finding_categories(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_audit_findings_severity_levels FOREIGN KEY (severity_level_id) REFERENCES severity_levels(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_audit_findings_audit_finding_statuses FOREIGN KEY (audit_finding_status_id) REFERENCES audit_finding_statuses(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE audit_findings;
-- +goose StatementEnd
