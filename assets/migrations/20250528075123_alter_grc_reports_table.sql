-- +goose Up
-- +goose StatementBegin
ALTER TABLE grc_reports
ADD reviewed_at DATETIME NULL;

ALTER TABLE grc_reports
ADD approved_at DATETIME NULL;

ALTER TABLE grc_reports
ADD grc_report_workflow_status_id INT NULL;

ALTER TABLE grc_reports
ADD CONSTRAINT FK_grc_reports_grc_report_workflow_statuses
FOREIGN KEY (grc_report_workflow_status_id) REFERENCES grc_report_workflow_statuses(id)
ON DELETE NO ACTION ON UPDATE NO ACTION;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE grc_reports
DROP CONSTRAINT FK_grc_reports_grc_report_workflow_statuses;

ALTER TABLE grc_reports
DROP COLUMN grc_report_workflow_status_id;

ALTER TABLE grc_reports
DROP COLUMN reviewed_at;

ALTER TABLE grc_reports
DROP COLUMN approved_at;
-- +goose StatementEnd
