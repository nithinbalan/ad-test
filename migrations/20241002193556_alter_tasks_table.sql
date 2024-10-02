-- +goose Up
-- +goose StatementBegin
ALTER TABLE tasks
ADD document_id INT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE tasks
DROP COLUMN document_id;
-- +goose StatementEnd
