-- +goose Up
-- +goose StatementBegin

UPDATE compliance_register
SET is_requirement_assessment = 1
WHERE is_requirement_assessment = 0
  AND is_control_assessment = 0;


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

UPDATE compliance_register
SET is_requirement_assessment = 0
WHERE is_requirement_assessment = 1
  AND is_control_assessment = 0;

-- +goose StatementEnd
