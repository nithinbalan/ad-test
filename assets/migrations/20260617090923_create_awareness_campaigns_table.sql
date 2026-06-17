-- +goose Up
-- +goose StatementBegin
CREATE TABLE awareness_campaigns(
    id INT PRIMARY KEY IDENTITY(1,1),
    reference_number VARCHAR(50) NOT NULL UNIQUE,
    awareness_campaign_type_id INT NOT NULL,

    opened_mail_percentage INT NULL,
    clicked_phishing_link_percentage INT NULL,
    entered_data_percentage INT NULL,
    reported_email_percentage INT NULL,
    risk_rating VARCHAR(50) NULL,
    campaign_result VARCHAR(50) NULL,
    recommendations NVARCHAR(MAX) NULL,

    campaign_title NVARCHAR(255) NOT NULL,
    message NVARCHAR(MAX) NULL,

    priority VARCHAR(20) NOT NULL,
    CONSTRAINT CHK_awareness_campaigns_priority
        CHECK (priority IN ('high', 'medium', 'low')),

    audience_id INT NULL,
    organization_unit_id INT NULL,
    target_scope_id INT NULL,
    supporting_unit_id INT NULL,

    target_date DATETIME NULL,

    is_survey BIT NOT NULL DEFAULT 0,
    survey_link NVARCHAR(MAX) NULL,
    content_link NVARCHAR(MAX) NULL,

    awareness_campaign_status_id INT NOT NULL,

    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL,
    updated_by_id INT NULL,
    updated_at DATETIME NULL,

    CONSTRAINT FK_awareness_campaigns_type FOREIGN KEY (awareness_campaign_type_id) REFERENCES awareness_campaign_types(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_awareness_campaigns_users_1 FOREIGN KEY (audience_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_awareness_campaigns_organization_units_1 FOREIGN KEY (organization_unit_id) REFERENCES organization_units(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_awareness_campaigns_users_2 FOREIGN KEY (target_scope_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_awareness_campaigns_organization_units_2 FOREIGN KEY (supporting_unit_id) REFERENCES organization_units(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_awareness_campaigns_status FOREIGN KEY (awareness_campaign_status_id) REFERENCES awareness_campaign_statuses(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_awareness_campaigns_users_3 FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_awareness_campaigns_users_4 FOREIGN KEY (updated_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE awareness_campaigns;
-- +goose StatementEnd
