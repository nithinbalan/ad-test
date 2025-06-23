-- +goose Up
-- +goose StatementBegin
UPDATE task_statuses
SET color_code = '#28A745'
WHERE type = 'closed';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Revert the color_code back to its previous value (assuming it was '#FFC107')
UPDATE task_statuses
SET color_code = '#008000'
WHERE type = 'closed';
-- +goose StatementEnd