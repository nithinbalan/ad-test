-- +goose Up
-- +goose StatementBegin
ALTER TABLE user_job_description_certificates
ADD is_mandatory BIT NOT NULL DEFAULT 0;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE user_job_description_certificates
DROP COLUMN is_mandatory;
-- +goose StatementEnd
