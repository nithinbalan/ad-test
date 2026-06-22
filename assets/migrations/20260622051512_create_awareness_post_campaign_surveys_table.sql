-- +goose Up
-- +goose StatementBegin
CREATE TABLE awareness_post_campaign_surveys(
    id INT PRIMARY KEY IDENTITY(1,1),
    reference_number VARCHAR(50) NOT NULL UNIQUE,
    linked_campaign_id INT NULL,
    send_date DATETIME NULL,
    target_audience NVARCHAR(MAX) NULL,
    recipients INT NULL,
    response_rate INT NULL,
    question_count INT NULL,
    average_rating INT NULL,
    maturity_level_id INT NULL,
    recommendations NVARCHAR(MAX) NULL,
    executing_party NVARCHAR(MAX) NULL,
    results_link NVARCHAR(MAX) NULL,
    created_by_id INT NOT NULL,
    updated_by_id INT NULL, 
    created_at DATETIME NOT NULL, 
    updated_at DATETIME NULL

    CONSTRAINT FK_awareness_post_campaign_surveys_campaign FOREIGN KEY (linked_campaign_id)
      REFERENCES awareness_campaigns(id) ON DELETE NO ACTION ON UPDATE NO ACTION, 
    CONSTRAINT FK_awareness_post_campaign_surveys_maturity_level FOREIGN KEY (maturity_level_id)
      REFERENCES awareness_maturity_levels(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_awareness_post_campaign_surveys_users_1 FOREIGN KEY (created_by_id) 
      REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION, 
    CONSTRAINT FK_awareness_post_campaign_surveys_users_2 FOREIGN KEY (updated_by_id) 
      REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE awareness_post_campaign_surveys;
-- +goose StatementEnd
