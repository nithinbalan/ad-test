-- +goose Up
-- +goose StatementBegin
UPDATE document_statuses
SET color_code = '#28A745'
WHERE type = 'published';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
UPDATE document_statuses
SET color_code = '#6C757D' 
WHERE type = 'published';
-- +goose StatementEnd
