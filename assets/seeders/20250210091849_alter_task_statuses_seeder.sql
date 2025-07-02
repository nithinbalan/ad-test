-- +goose Up
-- +goose StatementBegin
UPDATE task_statuses
SET color_code = '#FFC107'
WHERE type = 'in_progress';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Revert the color_code back to its previous value (assuming it was '#FFFF00')
UPDATE task_statuses
SET color_code = '#FFFF00'
WHERE type = 'in_progress';
-- +goose StatementEnd