-- +goose Up
-- +goose StatementBegin
CREATE TABLE identity_review_action_plan_comments (
    id INT PRIMARY KEY IDENTITY(1, 1),
    identity_review_action_plan_id INT NOT NULL,
    comment NVARCHAR(MAX) NOT NULL,
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL,

    CONSTRAINT FK_identity_review_action_plan_comments_plan
        FOREIGN KEY (identity_review_action_plan_id)
        REFERENCES identity_review_action_plans(id)
        ON DELETE CASCADE ON UPDATE CASCADE,

    CONSTRAINT FK_identity_review_action_plan_comments_users_1
        FOREIGN KEY (created_by_id) REFERENCES users(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE identity_review_action_plan_comments;
-- +goose StatementEnd
