-- +goose Up
-- +goose StatementBegin
UPDATE workflow_statuses
SET color_code = '#007BFF'
WHERE type = 'submit';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
UPDATE workflow_statuses
SET color_code = '#90DD75'
WHERE type = 'submit';
-- +goose StatementEnd
