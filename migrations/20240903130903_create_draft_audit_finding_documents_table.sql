-- +goose Up
-- +goose StatementBegin
CREATE TABLE draft_audit_finding_documents
(
    id INT PRIMARY KEY IDENTITY(1,1),
    draft_audit_finding_id INT NOT NULL,
    url VARCHAR(1000) DEFAULT NULL,
    name NVARCHAR(255) NULL,
    size INT DEFAULT NULL,
    token VARCHAR(255) DEFAULT NULL,
    original_name NVARCHAR(255) NULL,
    type VARCHAR(255) DEFAULT NULL,
    created_by_id INT,
    updated_by_id INT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,
    CONSTRAINT FK_draft_audit_finding_documents FOREIGN KEY (draft_audit_finding_id) REFERENCES draft_audit_findings(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_draft_audit_finding_documents_users_1 FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_draft_audit_finding_documents_users_2 FOREIGN KEY (updated_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE draft_audit_finding_documents;
-- +goose StatementEnd
