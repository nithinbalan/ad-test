-- +goose Up
-- +goose StatementBegin
ALTER TABLE consultation_request_comments
ALTER COLUMN consultation_request_status_id INT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE consultation_request_comments
ALTER COLUMN consultation_request_status_id INT NULL;
-- +goose StatementEnd