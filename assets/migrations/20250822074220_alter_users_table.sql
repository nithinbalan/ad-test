-- +goose Up
-- +goose StatementBegin
ALTER TABLE users
ADD delegatee_user_id INT NULL;

ALTER TABLE users
ADD CONSTRAINT FK_users_delegatee_user
FOREIGN KEY (delegatee_user_id) REFERENCES users(id)
ON DELETE NO ACTION ON UPDATE NO ACTION;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users
DROP CONSTRAINT FK_users_delegatee_user;

ALTER TABLE users
DROP COLUMN delegatee_user_id;
-- +goose StatementEnd
