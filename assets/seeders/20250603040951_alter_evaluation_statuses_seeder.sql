-- +goose Up
-- +goose StatementBegin
UPDATE evaluation_statuses
SET color_code = '#DC3545'
WHERE type = 'not-met';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
UPDATE evaluation_statuses
SET color_code = '#007BFF'
WHERE type = 'not-met';
-- +goose StatementEnd
