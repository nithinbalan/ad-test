-- +goose Up
-- +goose StatementBegin
EXEC sp_rename 'document_versions.enable_acknowledge_comment', 'enable_acknowledgement', 'COLUMN';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
EXEC sp_rename 'document_versions.enable_acknowledgement', 'enable_acknowledge_comment', 'COLUMN';
-- +goose StatementEnd

