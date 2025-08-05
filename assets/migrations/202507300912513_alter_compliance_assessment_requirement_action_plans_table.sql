-- +goose Up
-- +goose StatementBegin
-- Step 1: Drop FK and column for requirement_id
ALTER TABLE compliance_assessment_requirement_action_plans
DROP CONSTRAINT FK_compliance_assessment_requirement_action_plans_compliance_assessment_requirements;

ALTER TABLE compliance_assessment_requirement_action_plans
DROP COLUMN compliance_assessment_requirement_id;

-- Step 2: Add new column and FK to responsible_unit table
ALTER TABLE compliance_assessment_requirement_action_plans
ADD compliance_assessment_requirement_responsible_unit_id INT NOT NULL;

ALTER TABLE compliance_assessment_requirement_action_plans
ADD CONSTRAINT FK_carr_action_plans_responsible_unit
FOREIGN KEY (compliance_assessment_requirement_responsible_unit_id)
REFERENCES compliance_assessment_requirement_responsible_units(id)
ON DELETE NO ACTION ON UPDATE NO ACTION;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Revert the above changes

ALTER TABLE compliance_assessment_requirement_action_plans
DROP CONSTRAINT FK_carr_action_plans_responsible_unit;

ALTER TABLE compliance_assessment_requirement_action_plans
DROP COLUMN compliance_assessment_requirement_responsible_unit_id;

ALTER TABLE compliance_assessment_requirement_action_plans
ADD compliance_assessment_requirement_id INT NOT NULL;

ALTER TABLE compliance_assessment_requirement_action_plans
ADD CONSTRAINT FK_compliance_assessment_requirement_action_plans_compliance_assessment_requirements
FOREIGN KEY (compliance_assessment_requirement_id)
REFERENCES compliance_assessment_requirements(id)
ON DELETE NO ACTION ON UPDATE NO ACTION;
-- +goose StatementEnd
