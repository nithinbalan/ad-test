-- +goose Up
ALTER TABLE compliance_action_plans
ADD compliance_assessment_control_id INT NULL;

ALTER TABLE compliance_action_plans
ADD CONSTRAINT FK_compliance_action_plans_compliance_assessment_controls
FOREIGN KEY (compliance_assessment_control_id)
REFERENCES compliance_assessment_controls(id)
ON UPDATE CASCADE
ON DELETE CASCADE;

-- +goose Down
ALTER TABLE compliance_action_plans
DROP CONSTRAINT FK_compliance_action_plans_compliance_assessment_controls;

ALTER TABLE compliance_action_plans
DROP COLUMN compliance_assessment_control_id;
