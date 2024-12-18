-- +goose Up
-- +goose StatementBegin
CREATE TABLE audit_report_edit_histories (
    id INT PRIMARY KEY IDENTITY(1,1),
    audit_id INT NOT NULL,
    audit_draft_report_id INT DEFAULT NULL,
    audit_final_report_id INT DEFAULT NULL,
    action_type VARCHAR(20) NOT NULL DEFAULT 'create',
     CONSTRAINT CHK_audit_report_edit_histories_action_type CHECK (action_type IN ('create', 'update', 'submit', 'revert', 'approve', 'publish')),
    comment NVARCHAR(MAX) DEFAULT NULL,
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL,

    CONSTRAINT FK_audit_report_edit_histories_audits FOREIGN KEY (audit_id) REFERENCES audits(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_audit_report_edit_histories_draft_reports FOREIGN KEY (audit_draft_report_id) REFERENCES audit_draft_reports(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_audit_report_edit_histories_final_reports FOREIGN KEY (audit_final_report_id) REFERENCES audit_final_reports(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_audit_report_edit_histories_users FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS audit_report_edit_histories;
-- +goose StatementEnd
