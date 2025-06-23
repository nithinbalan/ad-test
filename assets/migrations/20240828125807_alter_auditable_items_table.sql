-- +goose Up
-- +goose StatementBegin
ALTER TABLE auditable_items
ADD responsible_user_id INT NULL;

ALTER TABLE auditable_items
ADD CONSTRAINT FK_auditable_items_users_4 FOREIGN KEY (responsible_user_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE auditable_items
DROP CONSTRAINT FK_auditable_items_users_4;

ALTER TABLE auditable_items
DROP COLUMN responsible_user_id;
-- +goose StatementEnd
