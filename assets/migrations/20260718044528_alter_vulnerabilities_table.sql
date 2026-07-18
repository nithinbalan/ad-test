-- +goose Up
-- +goose StatementBegin
ALTER TABLE vulnerabilities
ADD
    submitted_by_id INT NULL,
    submitted_at DATETIME NULL,
    next_review_level INT NULL,

    CONSTRAINT FK_vulnerabilities_users_submitted_by
        FOREIGN KEY (submitted_by_id)
        REFERENCES users(id)
        ON DELETE NO ACTION
        ON UPDATE NO ACTION;
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
ALTER TABLE vulnerabilities
DROP CONSTRAINT FK_vulnerabilities_users_submitted_by;

ALTER TABLE vulnerabilities
DROP COLUMN submitted_by_id,
            submitted_at,
            next_review_level;
-- +goose StatementEnd