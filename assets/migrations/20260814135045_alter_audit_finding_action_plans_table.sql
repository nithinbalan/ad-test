-- +goose Up
-- +goose StatementBegin

ALTER TABLE audit_finding_action_plans
ADD audit_finding_reassessment_id INT NULL;

ALTER TABLE audit_finding_action_plans
ADD CONSTRAINT FK_audit_finding_action_plans_reassessment
    FOREIGN KEY (audit_finding_reassessment_id)
    REFERENCES audit_finding_reassessments(id);

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

ALTER TABLE audit_finding_action_plans
DROP CONSTRAINT FK_audit_finding_action_plans_reassessment;

ALTER TABLE audit_finding_action_plans
DROP COLUMN audit_finding_reassessment_id;

-- +goose StatementEnd