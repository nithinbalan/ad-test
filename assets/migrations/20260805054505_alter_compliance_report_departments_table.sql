-- +goose Up
-- +goose StatementBegin
ALTER TABLE compliance_report_departments
ADD compliance_register_id INT NULL,
compliance_status_id INT NULL;

ALTER TABLE compliance_report_departments
ADD is_risk_based_assessment BIT NULL;

ALTER TABLE compliance_report_departments
ADD is_regulation_notification_assessment BIT NULL;

ALTER TABLE compliance_report_departments
ADD CONSTRAINT FK_compliance_report_departments_compliance_register FOREIGN KEY (compliance_register_id)
    REFERENCES compliance_register (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
CONSTRAINT FK_compliance_report_departments_compliance_statuses FOREIGN KEY (compliance_status_id)
    REFERENCES compliance_statuses(id) ON DELETE NO ACTION ON UPDATE NO ACTION;
;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE compliance_report_departments
DROP CONSTRAINT FK_compliance_report_departments_compliance_register;

ALTER TABLE compliance_report_departments
DROP CONSTRAINT FK_compliance_report_departments_compliance_statuses;

ALTER TABLE compliance_report_departments
DROP COLUMN is_regulation_notification_assessment;

ALTER TABLE compliance_report_departments
DROP COLUMN is_risk_based_assessment;

ALTER TABLE compliance_report_departments
DROP COLUMN compliance_register_id, compliance_status_id;
-- +goose StatementEnd