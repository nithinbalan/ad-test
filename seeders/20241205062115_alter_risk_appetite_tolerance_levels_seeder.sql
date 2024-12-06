-- +goose Up
-- +goose StatementBegin
UPDATE risk_appetite_tolerance_levels
SET color_code = '#28A745'
WHERE type = 'below';

UPDATE risk_appetite_tolerance_levels
SET color_code = '#FFC107'
WHERE type = 'within';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Reverting to the original color codes (if needed)
UPDATE risk_appetite_tolerance_levels
SET color_code = '#00FF00'
WHERE type = 'below';

UPDATE risk_appetite_tolerance_levels
SET color_code = '#FFFF00'
WHERE type = 'within';
-- +goose StatementEnd
