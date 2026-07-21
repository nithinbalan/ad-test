-- +goose Up
-- +goose StatementBegin
ALTER TABLE vulnerabilities ADD rejected_by_id INT NULL;
ALTER TABLE vulnerabilities ADD rejected_at DATETIME NULL;

ALTER TABLE vulnerabilities ADD CONSTRAINT FK_vulnerabilities_users_6
    FOREIGN KEY (rejected_by_id) REFERENCES users(id) ON DELETE NO ACTION;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE vulnerabilities DROP CONSTRAINT FK_vulnerabilities_users_6;
ALTER TABLE vulnerabilities DROP COLUMN rejected_by_id;
ALTER TABLE vulnerabilities DROP COLUMN rejected_at;
-- +goose StatementEnd
