-- +goose Up
-- +goose StatementBegin
ALTER TABLE compliance_assessments
ALTER COLUMN compliance_register_id INT NULL;

ALTER TABLE compliance_assessments
ADD type VARCHAR(20) NULL CONSTRAINT DF_compliance_assessments_type DEFAULT 'framework',
    responsible_unit_id INT NULL,
    CONSTRAINT CHK_compliance_assessments_type CHECK (type IN ('framework', 'department')),
    CONSTRAINT FK_compliance_assessments_responsible_unit FOREIGN KEY (responsible_unit_id) REFERENCES organization_units(id) ON DELETE NO ACTION ON UPDATE NO ACTION;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE compliance_assessments
DROP CONSTRAINT FK_compliance_assessments_responsible_unit;

ALTER TABLE compliance_assessments
DROP CONSTRAINT CHK_compliance_assessments_type;

ALTER TABLE compliance_assessments
DROP CONSTRAINT DF_compliance_assessments_type;

ALTER TABLE compliance_assessments
DROP COLUMN responsible_unit_id;

ALTER TABLE compliance_assessments
DROP COLUMN type;

ALTER TABLE compliance_assessments
ALTER COLUMN compliance_register_id INT NULL;
-- +goose StatementEnd
