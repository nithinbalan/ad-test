-- +goose Up
-- +goose StatementBegin
ALTER TABLE compliance_action_plans
ADD marked_as_auditable_item BIT NOT NULL 
    CONSTRAINT DF_compliance_action_plans_marked_as_auditable_item DEFAULT 0;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE compliance_action_plans
DROP CONSTRAINT DF_compliance_action_plans_marked_as_auditable_item;

ALTER TABLE compliance_action_plans
DROP COLUMN marked_as_auditable_item;
-- +goose StatementEnd
