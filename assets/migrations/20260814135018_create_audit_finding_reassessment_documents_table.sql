-- +goose Up
-- +goose StatementBegin

CREATE TABLE audit_finding_reassessment_documents(
    id INT PRIMARY KEY IDENTITY(1,1),
    audit_finding_reassessment_id INT NOT NULL,
    url VARCHAR(1000) DEFAULT NULL,
    name NVARCHAR(255) NULL,
    size INT DEFAULT NULL,
    token VARCHAR(255) DEFAULT NULL,
    original_name NVARCHAR(255) NULL,
    type VARCHAR(255) DEFAULT NULL,
    created_by_id INT NOT NULL,
    updated_by_id INT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,

    CONSTRAINT FK_audit_finding_reassessment_documents
        FOREIGN KEY (audit_finding_reassessment_id)
        REFERENCES audit_finding_reassessments(id),

    CONSTRAINT FK_audit_finding_reassessment_documents_created_by
        FOREIGN KEY (created_by_id)
        REFERENCES users(id),

    CONSTRAINT FK_audit_finding_reassessment_documents_updated_by
        FOREIGN KEY (updated_by_id)
        REFERENCES users(id)
);

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

DROP TABLE audit_finding_reassessment_documents;

-- +goose StatementEnd