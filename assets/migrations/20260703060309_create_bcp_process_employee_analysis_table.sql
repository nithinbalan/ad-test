-- +goose Up
-- +goose StatementBegin

CREATE TABLE bcp_process_employee_analysis (
    id INT PRIMARY KEY IDENTITY(1, 1),
    bcp_process_id INT NOT NULL,
    designation_id INT NOT NULL,
    normal_count INT NULL,
    critical_count INT NULL,
    created_by_id INT,
    created_at DATETIME NOT NULL,

     CONSTRAINT FK_bcp_process_employee_analysis_process FOREIGN KEY (bcp_process_id) REFERENCES bcp_processes(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_bcp_process_employee_analysis_designations FOREIGN KEY (designation_id) REFERENCES designations(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_bcp_process_employee_analysis_users_1 FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

DROP TABLE bcp_process_employee_analysis;

-- +goose StatementEnd
