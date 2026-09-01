-- +goose Up
-- +goose StatementBegin
ALTER TABLE auditable_items
ADD compliance_action_plan_id INT NULL;

ALTER TABLE auditable_items
ADD CONSTRAINT FK_auditable_items_compliance_action_plans
    FOREIGN KEY (compliance_action_plan_id) REFERENCES compliance_action_plans(id)
    ON DELETE NO ACTION ON UPDATE NO ACTION;

ALTER TABLE auditable_items
DROP CONSTRAINT CHK_auditable_items_type;

ALTER TABLE auditable_items
ADD CONSTRAINT CHK_auditable_items_type CHECK (type IN ('process', 'control', 'incident', 'risk',
 'document', 'compliance_register', 'compliance_requirement', 'compliance_action_plan', 'user', 'custom'));
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE auditable_items
DROP CONSTRAINT CHK_auditable_items_type;

ALTER TABLE auditable_items
ADD CONSTRAINT CHK_auditable_items_type CHECK (type IN ('process', 'control', 'incident', 'risk',
 'document', 'compliance_register', 'compliance_requirement', 'user', 'custom'));

ALTER TABLE auditable_items
DROP CONSTRAINT FK_auditable_items_compliance_action_plans;

ALTER TABLE auditable_items
DROP COLUMN compliance_action_plan_id;
-- +goose StatementEnd
