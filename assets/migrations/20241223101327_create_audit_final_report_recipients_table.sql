-- +goose Up
-- +goose StatementBegin
-- Create the audit_final_report_recipients table
CREATE TABLE audit_final_report_recipients (
    id INT PRIMARY KEY IDENTITY(1,1),
    audit_final_report_id INT NOT NULL, 
    user_id INT NOT NULL, 
    type VARCHAR(20) NOT NULL DEFAULT 'cc',
    CONSTRAINT CHK_audit_final_report_recipients_type CHECK (type IN ('cc', 'to')),
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL,

    -- Foreign key constraints
    CONSTRAINT FK_audit_final_report_recipients_audit_final_reports FOREIGN KEY (audit_final_report_id) REFERENCES audit_final_reports(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT FK_audit_final_report_recipients_users_1 FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT FK_audit_final_report_recipients_users_2 FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION

);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Drop the audit_final_report_recipients table
DROP TABLE IF EXISTS audit_final_report_recipients;
-- +goose StatementEnd
