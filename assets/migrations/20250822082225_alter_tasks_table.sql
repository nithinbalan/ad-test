-- +goose Up
-- +goose StatementBegin
ALTER TABLE tasks
ADD delegatee_user_id INT NULL;

ALTER TABLE tasks
ADD CONSTRAINT FK_tasks_delegatee_user
FOREIGN KEY (delegatee_user_id) REFERENCES users(id)
ON DELETE NO ACTION ON UPDATE NO ACTION;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE tasks
DROP CONSTRAINT FK_tasks_delegatee_user;

ALTER TABLE tasks
DROP COLUMN delegatee_user_id;
-- +goose StatementEnd
