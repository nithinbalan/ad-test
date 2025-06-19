-- +goose Up
-- +goose StatementBegin
ALTER TABLE user_job_description_competencies
ADD expected_level_id INT NULL;

ALTER TABLE user_job_description_competencies
ADD CONSTRAINT FK_user_jdc_expected_level
FOREIGN KEY (expected_level_id) REFERENCES competency_levels(id) ON DELETE NO ACTION;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE user_job_description_competencies
DROP CONSTRAINT FK_user_jdc_expected_level;

ALTER TABLE user_job_description_competencies
DROP COLUMN expected_level_id;
-- +goose StatementEnd
