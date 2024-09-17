-- +goose Up
-- +goose StatementBegin
ALTER TABLE audit_draft_reports
ADD submitted_by_id INT NULL;

ALTER TABLE audit_draft_reports
ADD submitted_at DATETIME NULL;

ALTER TABLE audit_draft_reports
ADD audit_draft_report_status_id INT NULL;

ALTER TABLE audit_draft_reports
ADD CONSTRAINT FK_audit_draft_reports_users_3 FOREIGN KEY (submitted_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION;

ALTER TABLE audit_draft_reports
ADD CONSTRAINT FK_audit_draft_reports_audit_draft_report_statuses FOREIGN KEY (audit_draft_report_status_id) REFERENCES audit_draft_report_statuses(id) ON DELETE NO ACTION ON UPDATE NO ACTION;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE audit_draft_reports
DROP CONSTRAINT FK_audit_draft_reports_users_3;

ALTER TABLE audit_draft_reports
DROP CONSTRAINT FK_audit_draft_reports_audit_draft_report_statuses;

ALTER TABLE audit_draft_reports
DROP COLUMN submitted_by_id;

ALTER TABLE audit_draft_reports
DROP COLUMN submitted_at;

ALTER TABLE audit_draft_reports
DROP COLUMN audit_draft_report_status_id;
-- +goose StatementEnd
