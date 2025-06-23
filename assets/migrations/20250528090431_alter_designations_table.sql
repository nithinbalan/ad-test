-- +goose Up
-- +goose StatementBegin
ALTER TABLE designations
ADD job_description_id INT NULL;

ALTER TABLE designations
ADD CONSTRAINT FK_designations_job_descriptions
FOREIGN KEY (job_description_id) REFERENCES job_descriptions(id)
ON DELETE NO ACTION ON UPDATE NO ACTION;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE designations
DROP CONSTRAINT FK_designations_job_descriptions;

ALTER TABLE designations
DROP COLUMN job_description_id;
-- +goose StatementEnd
