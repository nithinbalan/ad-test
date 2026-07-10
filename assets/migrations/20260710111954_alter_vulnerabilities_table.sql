-- +goose Up
-- +goose StatementBegin
ALTER TABLE vulnerabilities ADD closed_by_id INT NULL;
ALTER TABLE vulnerabilities ADD closed_at DATETIME NULL;

ALTER TABLE vulnerabilities ADD CONSTRAINT FK_vulnerabilities_users_4
    FOREIGN KEY (closed_by_id) REFERENCES users(id) ON DELETE NO ACTION;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE vulnerabilities DROP CONSTRAINT FK_vulnerabilities_users_4;
ALTER TABLE vulnerabilities DROP COLUMN closed_by_id;
ALTER TABLE vulnerabilities DROP COLUMN closed_at;
-- +goose StatementEnd
