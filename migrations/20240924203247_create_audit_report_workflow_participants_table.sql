-- +goose Up
-- +goose StatementBegin
CREATE TABLE audit_report_workflow_participants (
    id INT PRIMARY KEY IDENTITY(1,1),
    audit_report_workflow_id INT NOT NULL,
    user_id INT NOT NULL,
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL,
    CONSTRAINT FK_audit_report_workflow_participants_audit_report_workflows FOREIGN KEY (audit_report_workflow_id) REFERENCES audit_report_workflows(id),
    CONSTRAINT FK_audit_report_workflow_participants_users_1 FOREIGN KEY (user_id) REFERENCES users(id),
    CONSTRAINT FK_audit_report_workflow_participants_users_2 FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.audit_report_workflow_participants;
-- +goose StatementEnd
