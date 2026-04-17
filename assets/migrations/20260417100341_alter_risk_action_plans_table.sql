-- +goose Up
-- +goose StatementBegin

ALTER TABLE risk_action_plans
ADD 
    responsible_unit_id INT NULL,
    cost INT NULL,
    start_date DATETIME NULL,
    notes NVARCHAR(MAX) NULL;

-- Foreign Keys
ALTER TABLE risk_action_plans
ADD CONSTRAINT FK_risk_action_plans_organization_units
FOREIGN KEY (responsible_unit_id)
REFERENCES organization_units(id)
ON DELETE NO ACTION ON UPDATE NO ACTION;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Drop constraints first
ALTER TABLE risk_action_plans DROP CONSTRAINT FK_risk_action_plans_organization_units;

-- Drop columns
ALTER TABLE risk_action_plans
DROP COLUMN 
    responsible_unit_id,
    cost,
    start_date,
    notes;

-- +goose StatementEnd
