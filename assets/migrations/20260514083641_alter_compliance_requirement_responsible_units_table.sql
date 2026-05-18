-- +goose Up
-- +goose StatementBegin

ALTER TABLE compliance_requirement_responsible_units
ADD 
    is_regulation_notification_assessment BIT NULL,
    is_risk_based_assessment BIT NULL;

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

ALTER TABLE compliance_requirement_responsible_units
DROP COLUMN 
    is_regulation_notification_assessment,
    is_risk_based_assessment;

-- +goose StatementEnd