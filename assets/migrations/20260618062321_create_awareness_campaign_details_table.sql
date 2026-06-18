-- +goose Up
-- +goose StatementBegin
CREATE TABLE awareness_campaign_details (
    id INT PRIMARY KEY IDENTITY(1,1),
    awareness_campaign_id INT NOT NULL,
    awareness_pack_name NVARCHAR(255) NULL,
    recipient_count INT NULL,
    sent_count INT NULL,
    email_open_rate INT NULL,
    click_rate INT NULL,
    data_submission_rate INT NULL,
    report_rate INT NULL,
    risk_classification_id INT NULL,
    campaign_result NVARCHAR(MAX) NULL,
    improvement_recommendations NVARCHAR(MAX) NULL,
    subtype_role VARCHAR(100) NULL,
    presenter NVARCHAR(255) NULL,
    agenda NVARCHAR(MAX) NULL,
    session_start_time TIME NULL,
    session_end_time TIME NULL,
    venue NVARCHAR(255) NULL,
    attendee_count INT NULL,
    awareness_evaluation_rating_id INT NULL,
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL,
    updated_by_id INT NULL,
    updated_at DATETIME NULL,

    CONSTRAINT FK_awareness_campaign_details_campaign FOREIGN KEY (awareness_campaign_id)
        REFERENCES awareness_campaigns(id) ON DELETE CASCADE ON UPDATE NO ACTION,
    CONSTRAINT FK_awareness_campaign_details_awareness_priorities FOREIGN KEY (risk_classification_id)
       REFERENCES awareness_priorities(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_awareness_campaign_details_evaluation_rating FOREIGN KEY (awareness_evaluation_rating_id)
        REFERENCES awareness_evaluation_ratings(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_awareness_campaign_details_users_1 FOREIGN KEY (created_by_id) 
        REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_awareness_campaign_details_users_2 FOREIGN KEY (updated_by_id) 
        REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE awareness_campaign_details;
-- +goose StatementEnd
