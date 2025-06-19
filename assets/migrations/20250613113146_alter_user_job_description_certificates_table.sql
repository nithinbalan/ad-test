-- +goose Up
-- +goose StatementBegin
-- Drop the foreign key constraint first
ALTER TABLE user_job_description_certificates
DROP CONSTRAINT FK_user_jd_jdcc;

-- Drop the column
ALTER TABLE user_job_description_certificates
DROP COLUMN job_description_competency_certification_id;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Re-add the column
ALTER TABLE user_job_description_certificates
ADD job_description_competency_certification_id INT NULL;

-- Re-add the constraint
ALTER TABLE user_job_description_certificates
ADD CONSTRAINT FK_user_jd_jdcc FOREIGN KEY (job_description_competency_certification_id)
REFERENCES job_description_competency_certifications(id) ON DELETE CASCADE;
-- +goose StatementEnd
