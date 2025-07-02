-- +goose Up
-- +goose StatementBegin
CREATE TABLE audit_draft_report_comments
(
    id INT PRIMARY KEY IDENTITY(1,1),
    audit_draft_report_id INT NOT NULL,
    audit_draft_finding_id INT NULL,
    audit_draft_report_comment_id INT NULL,
    type VARCHAR(50) NOT NULL,
    CONSTRAINT CHK_audit_draft_report_comments_type CHECK (type IN ('executive-summary', 'finding', 'conclusion')),
    comment NVARCHAR(MAX) NULL,
    is_resolved BIT NOT NULL DEFAULT 0,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,
    created_by_id INT NOT NULL,
    updated_by_id INT NULL,

    CONSTRAINT FK_audit_draft_report_comments_audit_draft_reports FOREIGN KEY (audit_draft_report_id) REFERENCES audit_draft_reports(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_audit_draft_report_comments_self FOREIGN KEY (audit_draft_report_comment_id) REFERENCES audit_draft_report_comments(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_audit_draft_report_comments_users_1 FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_audit_draft_report_comments_users_2 FOREIGN KEY (updated_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.audit_draft_report_comments;
-- +goose StatementEnd
