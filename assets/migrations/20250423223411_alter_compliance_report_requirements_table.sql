-- +goose Up
-- +goose StatementBegin
ALTER TABLE compliance_report_requirements
ADD average_score DECIMAL(10,2) NULL;

ALTER TABLE compliance_report_requirements
ADD average_compliance_status_id INT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE compliance_report_requirements
DROP COLUMN average_score;

ALTER TABLE compliance_report_requirements
DROP COLUMN average_compliance_status_id;
-- +goose StatementEnd
