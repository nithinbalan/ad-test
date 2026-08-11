-- +goose Up
-- +goose StatementBegin
CREATE TABLE identity_review_user_account_decision_histories (
    id INT PRIMARY KEY IDENTITY(1, 1),
    identity_review_user_id INT NOT NULL,
    account_decision_id INT NOT NULL,
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL,

    CONSTRAINT FK_identity_review_user_account_decision_histories_review_users
        FOREIGN KEY (identity_review_user_id) REFERENCES identity_review_users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_identity_review_user_account_decision_histories_account_decisions
        FOREIGN KEY (account_decision_id) REFERENCES account_decisions(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_identity_review_user_account_decision_histories_created_by
        FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE identity_review_user_account_decision_histories;
-- +goose StatementEnd
