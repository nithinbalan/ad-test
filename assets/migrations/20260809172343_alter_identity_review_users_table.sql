-- +goose Up
-- +goose StatementBegin
ALTER TABLE identity_review_users
ADD user_identifier NVARCHAR(255) NULL,
    last_login NVARCHAR(255) NULL,
    account_decision_id INT NULL,
    current_role NVARCHAR(MAX) NULL,
    review_comment NVARCHAR(MAX) NULL;

ALTER TABLE identity_review_users
ADD CONSTRAINT FK_identity_review_users_account_decisions
FOREIGN KEY (account_decision_id) REFERENCES account_decisions(id)
ON DELETE NO ACTION
ON UPDATE NO ACTION;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE identity_review_users
DROP CONSTRAINT FK_identity_review_users_account_decisions;

ALTER TABLE identity_review_users
DROP COLUMN user_identifier,
            last_login,
            account_decision_id,
            current_role,
            review_comment;
-- +goose StatementEnd
