-- +goose Up
-- +goose StatementBegin
CREATE TABLE job_description_competency_certifications (
    id INT PRIMARY KEY IDENTITY(1, 1),
    job_description_competency_id INT NOT NULL,
    title NVARCHAR(255) NOT NULL,
    is_mandatory BIT NOT NULL DEFAULT 0,
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,


    CONSTRAINT FK_jdcc_certifications FOREIGN KEY (job_description_competency_id) REFERENCES job_description_competencies(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_jdcc_created_by FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
    
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE job_description_competency_certifications;
-- +goose StatementEnd
