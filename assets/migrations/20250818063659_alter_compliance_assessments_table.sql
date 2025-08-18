-- +goose Up
-- +goose StatementBegin
ALTER TABLE compliance_assessments
DROP CONSTRAINT FK_compliance_assessments_organization_units;

ALTER TABLE compliance_assessments
DROP COLUMN organization_unit_id;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE compliance_assessments
ADD organization_unit_id INT NULL;

ALTER TABLE compliance_assessments
ADD CONSTRAINT FK_compliance_assessments_organization_units
    FOREIGN KEY (organization_unit_id)
    REFERENCES organization_units(id)
    ON DELETE NO ACTION ON UPDATE NO ACTION;
-- +goose StatementEnd
