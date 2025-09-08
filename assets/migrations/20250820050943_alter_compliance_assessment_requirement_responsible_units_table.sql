-- +goose Up
-- +goose StatementBegin
ALTER TABLE compliance_assessment_requirement_responsible_units
ADD completed_at DATETIME NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE compliance_assessment_requirement_responsible_units
DROP COLUMN completed_at;
-- +goose StatementEnd
