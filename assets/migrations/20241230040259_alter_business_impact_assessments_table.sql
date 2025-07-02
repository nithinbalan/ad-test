-- +goose Up
-- +goose StatementBegin

-- Make the risk_level_id column nullable
ALTER TABLE business_impact_assessments
ALTER COLUMN risk_level_id INT NULL;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Revert the risk_level_id column to NOT NULL
ALTER TABLE business_impact_assessments
ALTER COLUMN risk_level_id INT NOT NULL;

-- +goose StatementEnd
