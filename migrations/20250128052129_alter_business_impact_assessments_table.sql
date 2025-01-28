-- +goose Up
-- +goose StatementBegin

-- adding is_available field
ALTER TABLE business_impact_assessments
ADD is_available TINYINT NOT NULL DEFAULT 0;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Revert the is_available
ALTER TABLE business_impact_assessments
DROP COLUMN is_available;

-- +goose StatementEnd
