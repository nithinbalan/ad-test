-- +goose Up
-- +goose StatementBegin
CREATE TABLE user_job_description_competencies (
    id INT PRIMARY KEY IDENTITY(1, 1),
    user_job_description_id INT NOT NULL,
    competency_id INT NOT NULL,
    actual_level_id INT NULL,
    comments NVARCHAR(MAX) NULL,
    competency_title_en NVARCHAR(255) NULL,
    competency_title_ar NVARCHAR(255) NULL,
    competency_group_title_en NVARCHAR(255) NULL,
    competency_group_title_ar NVARCHAR(255) NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT FK_user_jdc_user FOREIGN KEY (user_job_description_id) REFERENCES user_job_descriptions(id) ON DELETE CASCADE,
    CONSTRAINT FK_user_jdc_competencies FOREIGN KEY (competency_id) REFERENCES competencies(id) ON DELETE CASCADE,
    CONSTRAINT FK_user_jdc_level FOREIGN KEY (actual_level_id) REFERENCES competency_levels(id) ON DELETE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE user_job_description_competencies;
-- +goose StatementEnd