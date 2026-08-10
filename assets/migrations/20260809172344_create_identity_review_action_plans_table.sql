-- +goose Up
-- +goose StatementBegin
CREATE TABLE identity_review_action_plans (
    id INT PRIMARY KEY IDENTITY(1, 1),
    identity_review_user_id INT NULL,
    title NVARCHAR(500) NOT NULL,
    description NVARCHAR(MAX) NULL,
    severity_level_id INT NULL,
    target_date DATE NULL,
    responsible_user_id INT NOT NULL,
    review_user_id INT NULL,
    action_plan_status_id INT NOT NULL,
    percentage INT NOT NULL DEFAULT 0,
    comments NVARCHAR(MAX) NULL,
    delegatee_user_id INT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,
    created_by_id INT NOT NULL,
    updated_by_id INT NULL,

    CONSTRAINT FK_identity_review_action_plans_identity_review_users
        FOREIGN KEY (identity_review_user_id) REFERENCES identity_review_users(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_identity_review_action_plans_severity_levels
        FOREIGN KEY (severity_level_id) REFERENCES severity_levels(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_identity_review_action_plans_users_1
        FOREIGN KEY (responsible_user_id) REFERENCES users(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_identity_review_action_plans_users_4
        FOREIGN KEY (review_user_id) REFERENCES users(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_identity_review_action_plans_action_plan_statuses
        FOREIGN KEY (action_plan_status_id) REFERENCES action_plan_statuses(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_identity_review_action_plans_users_5
        FOREIGN KEY (delegatee_user_id) REFERENCES users(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_identity_review_action_plans_users_2
        FOREIGN KEY (created_by_id) REFERENCES users(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_identity_review_action_plans_users_3
        FOREIGN KEY (updated_by_id) REFERENCES users(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT CHK_identity_review_action_plans_percentage
        CHECK (percentage BETWEEN 0 AND 100)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE identity_review_action_plans;
-- +goose StatementEnd
