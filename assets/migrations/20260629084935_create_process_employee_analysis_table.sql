-- +goose Up
-- +goose StatementBegin

CREATE TABLE process_employee_analysis (
    id INT PRIMARY KEY IDENTITY(1, 1),
    process_id INT NOT NULL,
    designation_id INT NOT NULL,
    normal_count INT NULL,
    critical_count INT NULL,
    created_by_id INT,
    updated_by_id INT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,   

     CONSTRAINT FK_process_employee_analysis_process FOREIGN KEY (process_id) REFERENCES processes(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_process_employee_analysis_designations FOREIGN KEY (designation_id) REFERENCES designations(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_process_employee_analysis_users_1 FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_process_employee_analysis_users_2 FOREIGN KEY (updated_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

DROP TABLE process_employee_analysis;

-- +goose StatementEnd
