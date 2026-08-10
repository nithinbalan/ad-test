-- +goose Up
-- +goose StatementBegin
CREATE TABLE identity_review_campaign_workflow_histories (
    id INT PRIMARY KEY IDENTITY(1, 1),
    identity_review_campaign_id INT NOT NULL,
    level INT NOT NULL,
    workflow_status_id INT NOT NULL,
    comment NVARCHAR(MAX) NULL,
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT GETDATE(),

    CONSTRAINT FK_identity_review_campaign_workflow_histories_campaign
        FOREIGN KEY (identity_review_campaign_id) REFERENCES identity_review_campaigns(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_identity_review_campaign_workflow_histories_status
        FOREIGN KEY (workflow_status_id) REFERENCES workflow_statuses(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_identity_review_campaign_workflow_histories_users_1
        FOREIGN KEY (created_by_id) REFERENCES users(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE identity_review_campaign_workflow_histories;
-- +goose StatementEnd
