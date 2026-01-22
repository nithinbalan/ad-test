-- +goose Up
-- +goose StatementBegin
ALTER TABLE consultation_requests
ADD title NVARCHAR(255) NOT NULL;
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
ALTER TABLE consultation_requests
DROP COLUMN title;
-- +goose StatementEnd
