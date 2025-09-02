-- +goose Up
-- +goose StatementBegin

UPDATE document_statuses
SET color_code = '#FF9F45'
WHERE type = 'reverted';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

UPDATE document_statuses
SET color_code = '#28A745'
WHERE type = 'reverted';
-- +goose StatementEnd
