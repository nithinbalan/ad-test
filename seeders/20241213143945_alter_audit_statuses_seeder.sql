-- +goose Up
-- +goose StatementBegin
UPDATE audit_statuses
SET color_code = '#28A745'
WHERE type = 'completed';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Revert the color_code back to its previous value (assuming it was '#FFC107')
UPDATE audit_statuses
SET color_code = '#FFC107'
WHERE type = 'completed';
-- +goose StatementEnd
