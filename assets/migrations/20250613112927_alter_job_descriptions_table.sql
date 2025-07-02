-- +goose Up
-- +goose StatementBegin
ALTER TABLE job_descriptions
ADD created_at DATETIME NOT NULL DEFAULT GETDATE(),
    created_by_id INT NULL,
    updated_at DATETIME NULL,
    updated_by_id INT NULL;

ALTER TABLE job_descriptions
ADD CONSTRAINT FK_job_descriptions_users_created FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_job_descriptions_users_updated FOREIGN KEY (updated_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE job_descriptions DROP CONSTRAINT FK_job_descriptions_users_created;
ALTER TABLE job_descriptions DROP CONSTRAINT FK_job_descriptions_users_updated;

ALTER TABLE job_descriptions
DROP COLUMN created_at,
            created_by_id,
            updated_at,
            updated_by_id;
-- +goose StatementEnd
