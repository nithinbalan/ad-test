-- +goose Up
-- +goose StatementBegin
CREATE TABLE awareness_post_campaign_survey_themes ( 
    id INT PRIMARY KEY IDENTITY(1,1), 

    awareness_post_campaign_survey_id INT NOT NULL, 
    awareness_assessment_theme_id INT NOT NULL, 
    
    created_by_id INT NOT NULL, 
    created_at DATETIME NOT NULL, 

    CONSTRAINT FK_awareness_post_campaign_survey_themes_survey FOREIGN KEY (awareness_post_campaign_survey_id) 
      REFERENCES awareness_post_campaign_surveys(id) ON DELETE CASCADE ON UPDATE NO ACTION, 
    CONSTRAINT FK_awareness_post_campaign_survey_themes_theme FOREIGN KEY (awareness_assessment_theme_id) 
      REFERENCES awareness_assessment_themes(id) ON DELETE NO ACTION ON UPDATE NO ACTION, 
    CONSTRAINT FK_awareness_post_campaign_survey_themes_user_1 FOREIGN KEY (created_by_id) 
      REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION );
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE awareness_post_campaign_survey_themes;
-- +goose StatementEnd
