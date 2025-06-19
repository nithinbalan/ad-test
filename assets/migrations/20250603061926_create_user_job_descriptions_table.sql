-- +goose Up
-- +goose StatementBegin
CREATE TABLE user_job_descriptions (
    id INT PRIMARY KEY IDENTITY(1, 1),
    user_id INT NOT NULL,
    job_description_id INT NOT NULL,
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT FK_user_jds_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT FK_user_jds_jd FOREIGN KEY (job_description_id) REFERENCES job_descriptions(id) ON DELETE CASCADE,
    CONSTRAINT FK_user_jds_created_by FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION

);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE user_job_descriptions;
-- +goose StatementEnd