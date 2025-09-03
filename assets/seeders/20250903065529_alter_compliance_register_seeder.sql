-- +goose Up
-- +goose StatementBegin

UPDATE compliance_register
SET is_requirement_assessment = 1
WHERE is_requirement_assessment IS NULL
  AND is_control_assessment IS NULL;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

UPDATE compliance_register
SET is_requirement_assessment = NULL
WHERE is_requirement_assessment = 1
  AND is_control_assessment IS NULL;

-- +goose StatementEnd
