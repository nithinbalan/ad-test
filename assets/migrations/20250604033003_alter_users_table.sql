-- +goose Up
-- +goose StatementBegin
ALTER TABLE users
ADD reporting_user_id INT NULL;

ALTER TABLE users
ADD CONSTRAINT FK_users_job_descriptions
FOREIGN KEY (reporting_user_id) REFERENCES users(id)
ON DELETE NO ACTION ON UPDATE NO ACTION;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users
DROP CONSTRAINT FK_users_job_descriptions;

ALTER TABLE users
DROP COLUMN reporting_user_id;
-- +goose StatementEnd
