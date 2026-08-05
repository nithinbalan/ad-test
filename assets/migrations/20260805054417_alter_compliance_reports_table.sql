-- +goose Up
-- +goose StatementBegin
ALTER TABLE compliance_reports
ALTER COLUMN compliance_register_id INT NULL;

ALTER TABLE compliance_reports
ADD assessment_type NVARCHAR(20) NULL,
responsible_unit_id INT NULL;

ALTER TABLE compliance_reports
ADD CONSTRAINT FK_compliance_reports_organization_units FOREIGN KEY (responsible_unit_id)
REFERENCES organization_units (id) ON DELETE NO ACTION ON UPDATE NO ACTION
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE compliance_reports
ALTER COLUMN compliance_register_id INT NULL;

ALTER TABLE compliance_reports
DROP CONSTRAINT FK_compliance_reports_organization_units;

ALTER TABLE compliance_reports
DROP COLUMN assessment_type, responsible_unit_id;
-- +goose StatementEnd