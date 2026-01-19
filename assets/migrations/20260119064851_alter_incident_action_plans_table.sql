-- +goose Up
-- +goose StatementBegin

ALTER TABLE incident_action_plans
ADD action_plan_master_id INT NULL;

ALTER TABLE incident_action_plans
ADD CONSTRAINT FK_incident_action_plans_action_plan_master
    FOREIGN KEY (action_plan_master_id)
    REFERENCES action_plan_master(id)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE incident_action_plans
DROP CONSTRAINT FK_incident_action_plans_action_plan_master;

ALTER TABLE incident_action_plans
DROP COLUMN action_plan_master_id;

-- +goose StatementEnd
