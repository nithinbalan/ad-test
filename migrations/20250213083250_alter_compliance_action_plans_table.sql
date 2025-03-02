-- +goose Up
-- +goose StatementBegin
ALTER TABLE compliance_action_plans
ADD compliance_assessment_requirement_id INT NULL;

ALTER TABLE compliance_action_plans
ADD CONSTRAINT FK_compliance_action_plans_assessment_requirement FOREIGN KEY (compliance_assessment_requirement_id) REFERENCES compliance_assessment_requirements(id) ON DELETE NO ACTION ON UPDATE NO ACTION;

ALTER TABLE compliance_action_plans 
ALTER COLUMN compliance_register_id INT NULL;
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
ALTER TABLE compliance_action_plans
DROP CONSTRAINT FK_compliance_action_plans_assessment_requirement;

ALTER TABLE compliance_action_plans
DROP COLUMN compliance_assessment_requirement_id;

ALTER TABLE compliance_action_plans 
ALTER COLUMN compliance_register_id INT NOT NULL;
-- +goose StatementEnd
