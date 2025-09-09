-- +goose Up
-- +goose StatementBegin
ALTER TABLE document_acknowledge_users
ADD comment NVARCHAR(MAX) NULL;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE document_acknowledge_users
DROP COLUMN comment;
-- +goose StatementEnd

