-- +goose Up
-- +goose StatementBegin
UPDATE risk_levels
SET color_code = '#FFC107'
WHERE type = 'medium';

UPDATE risk_levels
SET color_code = '#28A745'
WHERE type = 'low';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
UPDATE risk_levels
SET color_code = '#FFFF00'
WHERE type = 'medium';

UPDATE risk_levels
SET color_code = '#00FF00'
WHERE type = 'low';
-- +goose StatementEnd
