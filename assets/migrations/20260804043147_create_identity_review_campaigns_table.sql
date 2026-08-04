-- +goose Up
-- +goose StatementBegin
CREATE TABLE identity_review_campaigns (
    id INT PRIMARY KEY IDENTITY(1, 1),
    reference_number VARCHAR(45) NOT NULL,
    campaign_name NVARCHAR(250) NOT NULL,
    review_type_id INT NOT NULL,
    description NVARCHAR(MAX) NULL,
    start_date DATE NOT NULL,
    due_date DATE NOT NULL,
    reviewer_id INT NOT NULL,
    approver_id INT NOT NULL,
    campaign_status_id INT NOT NULL,
    reminder_days_offset INT NULL,
    completion_percentage INT NOT NULL DEFAULT 0,
    created_by_id INT NOT NULL,
    updated_by_id INT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,

    CONSTRAINT FK_identity_review_campaigns_review_types
        FOREIGN KEY (review_type_id) REFERENCES review_types(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_identity_review_campaigns_users_1
        FOREIGN KEY (reviewer_id) REFERENCES users(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_identity_review_campaigns_users_2
        FOREIGN KEY (approver_id) REFERENCES users(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_identity_review_campaigns_review_campaign_statuses
        FOREIGN KEY (campaign_status_id) REFERENCES review_campaign_statuses(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_identity_review_campaigns_users_3
        FOREIGN KEY (created_by_id) REFERENCES users(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_identity_review_campaigns_users_4
        FOREIGN KEY (updated_by_id) REFERENCES users(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE identity_review_campaigns;
-- +goose StatementEnd