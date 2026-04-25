-- +goose Up
-- +goose StatementBegin
-- Update Partially Effective max_score to prevent overlap with Very Effective
-- This ensures scores of exactly 5.0 are classified as Very Effective only
UPDATE control_effectiveness_ratings
SET max_score = 4.99
WHERE type = 'partially_effective';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Revert Partially Effective max_score back to 5.0
UPDATE control_effectiveness_ratings
SET max_score = 5.0
WHERE type = 'partially_effective';
-- +goose StatementEnd
