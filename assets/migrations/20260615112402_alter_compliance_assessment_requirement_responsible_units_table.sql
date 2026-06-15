-- +goose Up
-- +goose StatementBegin

ALTER TABLE compliance_assessment_requirement_responsible_units
ADD 
    is_regulation_notification_assessment BIT NOT NULL
    CONSTRAINT DF_compliance_assessment_requirement_responsible_units_is_regulation_notification_assessment DEFAULT 0,
    
    is_risk_based_assessment BIT NOT NULL
    CONSTRAINT DF_compliance_assessment_requirement_responsible_units_is_risk_based_assessment DEFAULT 0;

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

ALTER TABLE compliance_assessment_requirement_responsible_units
DROP CONSTRAINT IF EXISTS DF_compliance_assessment_requirement_responsible_units_is_regulation_notification_assessment;

ALTER TABLE compliance_assessment_requirement_responsible_units
DROP CONSTRAINT IF EXISTS DF_compliance_assessment_requirement_responsible_units_is_risk_based_assessment;

ALTER TABLE compliance_assessment_requirement_responsible_units
DROP COLUMN 
    is_regulation_notification_assessment,
    is_risk_based_assessment;

-- +goose StatementEnd
