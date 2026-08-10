-- +goose Up
-- +goose StatementBegin
CREATE TABLE identity_review_campaign_review_comments (
    id INT PRIMARY KEY IDENTITY(1, 1),
    identity_review_campaign_id INT NOT NULL,
    comment NVARCHAR(1000) NOT NULL,
    comment_status_id INT NOT NULL,
    level INT NOT NULL,
    identity_review_campaign_review_comment_id INT NULL,
    created_by_id INT NOT NULL,
    updated_by_id INT NULL,
    created_at DATETIME NOT NULL DEFAULT GETDATE(),
    updated_at DATETIME NULL,

    CONSTRAINT FK_identity_review_campaign_review_comments_campaign
        FOREIGN KEY (identity_review_campaign_id) REFERENCES identity_review_campaigns(id)
        ON DELETE CASCADE ON UPDATE CASCADE,

    CONSTRAINT FK_identity_review_campaign_review_comments_status
        FOREIGN KEY (comment_status_id) REFERENCES comment_statuses(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_identity_review_campaign_review_comments_parent
        FOREIGN KEY (identity_review_campaign_review_comment_id) REFERENCES identity_review_campaign_review_comments(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_identity_review_campaign_review_comments_users_1
        FOREIGN KEY (created_by_id) REFERENCES users(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_identity_review_campaign_review_comments_users_2
        FOREIGN KEY (updated_by_id) REFERENCES users(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE identity_review_campaign_review_comments;
-- +goose StatementEnd
