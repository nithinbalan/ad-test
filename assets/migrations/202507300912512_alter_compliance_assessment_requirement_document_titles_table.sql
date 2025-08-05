-- +goose Up
-- +goose StatementBegin
-- Step 1: Drop old foreign key and column
ALTER TABLE compliance_assessment_requirement_document_titles
DROP CONSTRAINT FK_compliance_assessment_requirement_document_titles_compliance_assessment_requirements;

ALTER TABLE compliance_assessment_requirement_document_titles
DROP COLUMN compliance_assessment_requirement_id;

-- Step 2: Add new column and foreign key
ALTER TABLE compliance_assessment_requirement_document_titles
ADD compliance_assessment_requirement_responsible_unit_id INT NULL;

ALTER TABLE compliance_assessment_requirement_document_titles
ADD CONSTRAINT FK_document_titles_responsible_unit
FOREIGN KEY (compliance_assessment_requirement_responsible_unit_id)
REFERENCES compliance_assessment_requirement_responsible_units(id)
ON DELETE NO ACTION ON UPDATE NO ACTION;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Revert the above changes

ALTER TABLE compliance_assessment_requirement_document_titles
DROP CONSTRAINT FK_document_titles_responsible_unit;

ALTER TABLE compliance_assessment_requirement_document_titles
DROP COLUMN compliance_assessment_requirement_responsible_unit_id;

ALTER TABLE compliance_assessment_requirement_document_titles
ADD compliance_assessment_requirement_id INT NULL;

ALTER TABLE compliance_assessment_requirement_document_titles
ADD CONSTRAINT FK_compliance_assessment_requirement_document_titles_compliance_assessment_requirements
FOREIGN KEY (compliance_assessment_requirement_id)
REFERENCES compliance_assessment_requirements(id)
ON DELETE NO ACTION ON UPDATE NO ACTION;
-- +goose StatementEnd
