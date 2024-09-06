-- +goose Up
-- +goose StatementBegin
CREATE TABLE audit_draft_finding_conversations
(
    id INT PRIMARY KEY IDENTITY(1,1),
    audit_draft_finding_id INT NOT NULL,
    description NVARCHAR(MAX) NULL,
    type VARCHAR(20) NOT NULL,
    CONSTRAINT CHK_audit_draft_finding_conversations_type CHECK (type IN ('DR', 'AC')),
    created_by_id INT,
    updated_by_id INT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,
    CONSTRAINT FK_audit_draft_finding_conversations FOREIGN KEY (audit_draft_finding_id) REFERENCES audit_draft_findings(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_audit_draft_finding_conversations_users_1 FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_audit_draft_finding_conversations_users_2 FOREIGN KEY (updated_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE audit_draft_finding_conversations;
-- +goose StatementEnd
