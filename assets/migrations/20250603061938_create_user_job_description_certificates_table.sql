
-- +goose Up
-- +goose StatementBegin
CREATE TABLE user_job_description_certificates (
    id INT PRIMARY KEY IDENTITY(1,1),
    user_job_description_competency_id INT NOT NULL,
    job_description_competency_certification_id INT NOT NULL,
    certificate_title NVARCHAR(255) NOT NULL,
    url VARCHAR(1000) DEFAULT NULL,
    name NVARCHAR(255) NULL,
    size INT DEFAULT NULL,
    token VARCHAR(255) DEFAULT NULL,
    original_name NVARCHAR(255) NULL,
    type VARCHAR(255) DEFAULT NULL,
    created_by_id INT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT FK_user_jd_cert_user_jdc FOREIGN KEY (user_job_description_competency_id) REFERENCES user_job_description_competencies(id) ON DELETE CASCADE,
    CONSTRAINT FK_user_jd_jdcc FOREIGN KEY (job_description_competency_certification_id) REFERENCES job_description_competency_certifications(id) ON DELETE CASCADE,
    CONSTRAINT FK_user_jd_cert_created_by FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE user_job_description_certificates;
-- +goose StatementEnd