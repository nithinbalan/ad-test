-- +goose Up
-- +goose StatementBegin
CREATE TABLE job_description_competencies (
    id INT PRIMARY KEY IDENTITY(1, 1),
    job_description_id INT NOT NULL,
    competency_id INT NOT NULL,
    competency_level_id INT NULL,
    created_by_id INT NOT NULL,
    updated_by_id INT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME,

    CONSTRAINT FK_jdc_job_descriptions FOREIGN KEY (job_description_id) REFERENCES job_descriptions(id) ON DELETE CASCADE,
    CONSTRAINT FK_jdc_competencies FOREIGN KEY (competency_id) REFERENCES competencies(id) ON DELETE CASCADE,
    CONSTRAINT FK_jdc_competency_levels FOREIGN KEY (competency_level_id) REFERENCES competency_levels(id) ON DELETE NO ACTION,
    CONSTRAINT FK_jdc_created_by FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_jdc_updated_by FOREIGN KEY (updated_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE job_description_competencies;
-- +goose StatementEnd
