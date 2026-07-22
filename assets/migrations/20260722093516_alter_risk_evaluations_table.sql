-- +goose Up
-- +goose StatementBegin

ALTER TABLE risk_evaluations
ADD
    is_latest BIT NOT NULL
        CONSTRAINT DF_risk_evaluations_is_latest DEFAULT (1),
    risk_register_update_plan_id INT NULL;

ALTER TABLE risk_evaluations
ADD CONSTRAINT FK_risk_evaluations_risk_register_update_plans
FOREIGN KEY (risk_register_update_plan_id)
REFERENCES risk_register_update_plans(id)
ON DELETE NO ACTION
ON UPDATE NO ACTION;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE risk_evaluations
DROP CONSTRAINT FK_risk_evaluations_risk_register_update_plans;

ALTER TABLE risk_evaluations
DROP CONSTRAINT DF_risk_evaluations_is_latest;

ALTER TABLE risk_evaluations
DROP COLUMN
    is_latest,
    risk_register_update_plan_id;
-- +goose StatementEnd
