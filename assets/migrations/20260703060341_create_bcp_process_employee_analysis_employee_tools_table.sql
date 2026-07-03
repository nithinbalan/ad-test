-- +goose Up
-- +goose StatementBegin
CREATE TABLE bcp_process_employee_analysis_employee_tools (
    id INT IDENTITY(1,1) PRIMARY KEY,
    bcp_process_employee_analysis_employee_id INT NOT NULL,
    equipment_tool_id INT NOT NULL,
    is_has BIT NOT NULL DEFAULT 0,
    is_needed BIT NOT NULL DEFAULT 0,
    created_by_id INT,
    created_at DATETIME NOT NULL,
    CONSTRAINT FK_bcp_process_employee_analysis_employee_tools_bcp_process_employee_analysis_employee FOREIGN KEY (bcp_process_employee_analysis_employee_id) REFERENCES bcp_process_employee_analysis_employees(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_bcp_process_employee_analysis_employee_tools_equipment_tools FOREIGN KEY (equipment_tool_id) REFERENCES process_equipment_tools(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_bcp_process_employee_analysis_employee_tools_users_1 FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE bcp_process_employee_analysis_employee_tools;
-- +goose StatementEnd

