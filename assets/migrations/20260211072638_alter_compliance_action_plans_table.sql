-- +goose Up
-- +goose StatementBegin
ALTER TABLE compliance_action_plans
ADD compliance_assessment_impact_type_id INT NULL;

ALTER TABLE compliance_action_plans
ADD CONSTRAINT FK_compliance_action_plans_compliance_assessment_impact_type FOREIGN KEY (compliance_assessment_impact_type_id) REFERENCES compliance_assessment_impact_types(id) ON DELETE NO ACTION ON UPDATE NO ACTION;


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE compliance_action_plans
DROP CONSTRAINT FK_compliance_action_plans_compliance_assessment_impact_type;
ALTER TABLE compliance_action_plans
DROP COLUMN compliance_assessment_impact_type_id;
-- +goose StatementEnd
