-- +goose Up
-- +goose StatementBegin
UPDATE severity_levels 
SET color_code = '#FD7E14'
WHERE type = 'high';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
UPDATE severity_levels
SET color_code = '#DC3545'
WHERE type = 'high';
-- +goose StatementEnd
