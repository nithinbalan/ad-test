-- +goose Up
-- +goose StatementBegin
UPDATE document_change_request_statuses
SET color_code = '#28A745'
WHERE type = 'published';

UPDATE document_change_request_statuses
SET color_code = '#FF9F45'
WHERE type = 'reverted';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
UPDATE document_change_request_statuses
SET color_code = '#6C757D'
WHERE type = 'published';

UPDATE document_change_request_statuses
SET color_code = '#28A745'
WHERE type = 'reverted';
-- +goose StatementEnd
