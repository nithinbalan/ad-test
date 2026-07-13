-- +goose Up
-- +goose StatementBegin
ALTER TABLE vulnerabilities ADD deffered_by_id INT NULL;
ALTER TABLE vulnerabilities ADD deffered_at DATETIME NULL;

ALTER TABLE vulnerabilities ADD CONSTRAINT FK_vulnerabilities_users_5
    FOREIGN KEY (deffered_by_id) REFERENCES users(id) ON DELETE NO ACTION;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE vulnerabilities DROP CONSTRAINT FK_vulnerabilities_users_5;
ALTER TABLE vulnerabilities DROP COLUMN deffered_by_id;
ALTER TABLE vulnerabilities DROP COLUMN deffered_at;
-- +goose StatementEnd
