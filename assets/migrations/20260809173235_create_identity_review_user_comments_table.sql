-- +goose Up
-- +goose StatementBegin
CREATE TABLE identity_review_user_comments (
    id INT PRIMARY KEY IDENTITY(1, 1),
    identity_review_user_id INT NOT NULL,
    account_decision_id INT NULL,
    comment NVARCHAR(MAX) NOT NULL,
    created_by_id INT NOT NULL,
    updated_by_id INT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,

    CONSTRAINT FK_identity_review_user_comments_identity_review_users
        FOREIGN KEY (identity_review_user_id) REFERENCES identity_review_users(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_identity_review_user_comments_account_decisions
        FOREIGN KEY (account_decision_id) REFERENCES account_decisions(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_identity_review_user_comments_users_1
        FOREIGN KEY (created_by_id) REFERENCES users(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_identity_review_user_comments_users_2
        FOREIGN KEY (updated_by_id) REFERENCES users(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE identity_review_user_comments;
-- +goose StatementEnd
