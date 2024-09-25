-- +goose Up
-- +goose StatementBegin
CREATE TABLE audit_report_reviews (
    id INT PRIMARY KEY IDENTITY(1,1),
    audit_draft_report_id INT NOT NULL,
    audit_report_workflow_id INT NOT NULL,
    workflow_status_id INT NOT NULL,
    user_id INT NOT NULL,
    comment NVARCHAR(MAX) NULL,
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL,

    CONSTRAINT FK_audit_report_reviews_audit_draft_reports FOREIGN KEY (audit_draft_report_id) REFERENCES audit_draft_reports(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_audit_report_reviews_audit_report_workflows FOREIGN KEY (audit_report_workflow_id) REFERENCES audit_report_workflows(id),
    CONSTRAINT FK_audit_report_reviews_workflow_statuses FOREIGN KEY (workflow_status_id) REFERENCES workflow_statuses(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_audit_report_reviews_users_1 FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_audit_report_reviews_users_2 FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.audit_report_reviews;
-- +goose StatementEnd
