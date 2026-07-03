-- +goose Up
-- +goose StatementBegin
CREATE TABLE bcp_process_system_registers (
    id INT IDENTITY(1,1) PRIMARY KEY,
    bcp_process_id INT NOT NULL,
    system_register_id INT NOT NULL,
    created_by_id INT,
    created_at DATETIME NOT NULL,
    CONSTRAINT FK_bcp_process_system_registers_processes FOREIGN KEY (bcp_process_id) REFERENCES bcp_processes(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_bcp_process_system_registers_system_registers FOREIGN KEY (system_register_id) REFERENCES system_registers(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_bcp_process_system_registers_users_1 FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE bcp_process_system_registers;
-- +goose StatementEnd
