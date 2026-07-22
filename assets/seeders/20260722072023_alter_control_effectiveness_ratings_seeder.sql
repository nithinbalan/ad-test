-- +goose Up
-- +goose StatementBegin
UPDATE control_effectiveness_ratings
SET min_score = CASE type
    WHEN 'ineffective' THEN 0.11
    WHEN 'improvement_needed' THEN 2.01
    WHEN 'partially_effective' THEN 3.51
END
WHERE type IN ('ineffective', 'improvement_needed', 'partially_effective');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
UPDATE control_effectiveness_ratings
SET min_score = CASE type
    WHEN 'ineffective' THEN 0.1
    WHEN 'improvement_needed' THEN 2.0
    WHEN 'partially_effective' THEN 3.5
END
WHERE type IN ('ineffective', 'improvement_needed', 'partially_effective');
-- +goose StatementEnd
