-- +goose Up
-- +goose StatementBegin

CREATE TABLE audit_finding_reassessments(
    id INT PRIMARY KEY IDENTITY(1, 1),
    audit_finding_id INT NOT NULL,
    outcome VARCHAR(20) NULL,
    comment NVARCHAR(MAX) NULL,
    created_by_id INT NOT NULL,
    updated_by_id INT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,

    CONSTRAINT FK_audit_finding_reassessments_audit_findings
        FOREIGN KEY (audit_finding_id)
        REFERENCES audit_findings(id),

    CONSTRAINT FK_audit_finding_reassessments_created_by
        FOREIGN KEY (created_by_id)
        REFERENCES users(id),

    CONSTRAINT FK_audit_finding_reassessments_updated_by
        FOREIGN KEY (updated_by_id)
        REFERENCES users(id)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE audit_finding_reassessments;

-- +goose StatementEnd