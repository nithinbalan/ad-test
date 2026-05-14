-- +goose Up
-- +goose StatementBegin
CREATE TABLE system_register_job_titles (
    id INT PRIMARY KEY IDENTITY(1,1),
    system_register_id INT NOT NULL,
    designation_id INT NOT NULL,
    created_by_id INT,
    created_at DATETIME NOT NULL,
    CONSTRAINT FK_system_register_job_titles_system_registers FOREIGN KEY (system_register_id) REFERENCES system_registers(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_system_register_job_titles_designations FOREIGN KEY (designation_id) REFERENCES designations(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_system_register_job_titles_users_1 FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE system_register_job_titles;
-- +goose StatementEnd
