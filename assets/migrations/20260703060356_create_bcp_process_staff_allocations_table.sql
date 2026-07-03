-- +goose Up
-- +goose StatementBegin
CREATE TABLE bcp_process_staff_allocations (
    id INT IDENTITY(1,1) PRIMARY KEY,
    bcp_process_id INT NOT NULL,
    bcp_process_employee_analysis_id INT NULL,
    created_by_id INT,
    created_at DATETIME NOT NULL, 

    CONSTRAINT FK_bcp_process_staff_allocations_bcp_process FOREIGN KEY (bcp_process_id)
        REFERENCES bcp_processes(id) ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_bcp_process_staff_allocations_bcp_process_employee_analysis FOREIGN KEY (bcp_process_employee_analysis_id)
        REFERENCES bcp_process_employee_analysis(id) ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_bcp_process_staff_allocations_users_1 FOREIGN KEY (created_by_id)
        REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE bcp_process_staff_allocations;
-- +goose StatementEnd