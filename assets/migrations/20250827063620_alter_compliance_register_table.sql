-- +goose Up
-- +goose StatementBegin
ALTER TABLE compliance_register
ADD is_control_assessment BIT NOT NULL CONSTRAINT DF_compliance_register_is_control_assessment DEFAULT 0;

ALTER TABLE compliance_register
ADD is_requirement_assessment BIT NOT NULL CONSTRAINT DF_compliance_register_is_requirement_assessment DEFAULT 0;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE compliance_register
DROP CONSTRAINT DF_compliance_register_is_control_assessment;

ALTER TABLE compliance_register
DROP COLUMN is_control_assessment;

ALTER TABLE compliance_register
DROP CONSTRAINT DF_compliance_register_is_requirement_assessment;

ALTER TABLE compliance_register
DROP COLUMN is_requirement_assessment;
-- +goose StatementEnd
