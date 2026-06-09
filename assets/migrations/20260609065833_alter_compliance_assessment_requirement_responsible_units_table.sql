-- +goose Up
-- +goose StatementBegin
ALTER TABLE compliance_assessment_requirement_responsible_units
ADD is_not_applicable BIT NOT NULL  CONSTRAINT DF_compliance_requirement_responsible_is_not_applicable DEFAULT 0;

ALTER TABLE compliance_assessment_requirement_responsible_units
ADD compliance_status_id INT NULL;

ALTER TABLE compliance_assessment_requirement_responsible_units
ADD CONSTRAINT FK_compliance_assessment_requirement_responsible_units_compliance_status FOREIGN KEY (compliance_status_id) REFERENCES compliance_statuses(id) ON DELETE NO ACTION ON UPDATE NO ACTION;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE compliance_assessment_requirement_responsible_units
DROP CONSTRAINT DF_compliance_requirement_responsible_is_not_applicable;

ALTER TABLE compliance_assessment_requirement_responsible_units
DROP CONSTRAINT FK_compliance_assessment_requirement_responsible_units_compliance_status;

ALTER TABLE compliance_assessment_requirement_responsible_units
DROP COLUMN compliance_status_id;

ALTER TABLE compliance_assessment_requirement_responsible_units
DROP COLUMN is_not_applicable;
-- +goose StatementEnd
