-- +goose Up
-- +goose StatementBegin

ALTER TABLE compliance_requirements
ADD 
    is_regulation_notification_assessment BIT NULL,
    is_risk_based_assessment BIT NULL;

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

ALTER TABLE compliance_requirements
DROP COLUMN 
    is_regulation_notification_assessment,
    is_risk_based_assessment;

-- +goose StatementEnd