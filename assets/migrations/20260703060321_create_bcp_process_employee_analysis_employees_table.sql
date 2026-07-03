-- +goose Up
-- +goose StatementBegin
CREATE TABLE bcp_process_employee_analysis_employees (
    id INT IDENTITY(1,1) PRIMARY KEY,
    bcp_process_employee_analysis_id INT NOT NULL,
    user_id INT NOT NULL,
    created_by_id INT,
    created_at DATETIME NOT NULL,
    CONSTRAINT FK_bcp_process_employee_analysis_employees_bcp_process_employee_analysis FOREIGN KEY (bcp_process_employee_analysis_id) REFERENCES bcp_process_employee_analysis(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_bcp_process_employee_analysis_employees_users FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_bcp_process_employee_analysis_employees_users_1 FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE bcp_process_employee_analysis_employees;
-- +goose StatementEnd
