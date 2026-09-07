-- +goose Up
-- +goose StatementBegin
ALTER TABLE project_change_request_documents 
ALTER COLUMN token VARCHAR(MAX) NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE project_change_request_documents 
ALTER COLUMN token VARCHAR(100) NOT NULL;
-- +goose StatementEnd
