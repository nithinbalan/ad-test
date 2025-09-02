-- +goose Up
-- +goose StatementBegin
-- Step 1: Drop FK and column for requirement_id
ALTER TABLE compliance_action_plans
DROP CONSTRAINT FK_compliance_action_plans_assessment_requirement;

ALTER TABLE compliance_action_plans
DROP COLUMN compliance_assessment_requirement_id;

-- Step 2: Add new column and FK to responsible_unit table
ALTER TABLE compliance_action_plans
ADD compliance_assessment_requirement_responsible_unit_id INT NULL;

ALTER TABLE compliance_action_plans
ADD CONSTRAINT FK_compliance_action_plans_responsible_unit
FOREIGN KEY (compliance_assessment_requirement_responsible_unit_id)
REFERENCES compliance_assessment_requirement_responsible_units(id)
ON DELETE NO ACTION ON UPDATE NO ACTION;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Step 1: Drop new FK and column to responsible_unit table
ALTER TABLE compliance_action_plans
DROP CONSTRAINT FK_compliance_action_plans_responsible_unit;

ALTER TABLE compliance_action_plans
DROP COLUMN compliance_assessment_requirement_responsible_unit_id;

-- Step 2: Re-add old column and FK to assessment requirement
ALTER TABLE compliance_action_plans
ADD compliance_assessment_requirement_id INT NULL;

ALTER TABLE compliance_action_plans
ADD CONSTRAINT FK_compliance_action_plans_assessment_requirement
FOREIGN KEY (compliance_assessment_requirement_id)
REFERENCES compliance_assessment_requirements(id)
ON DELETE NO ACTION ON UPDATE NO ACTION;
-- +goose StatementEnd

