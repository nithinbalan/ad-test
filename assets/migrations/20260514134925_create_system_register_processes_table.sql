-- +goose Up
-- +goose StatementBegin
CREATE TABLE system_register_processes (
    id INT PRIMARY KEY IDENTITY(1, 1),
    system_register_id INT NOT NULL,
    process_id INT NOT NULL,
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL,
    CONSTRAINT FK_system_register_processes_system_register FOREIGN KEY (system_register_id) REFERENCES system_registers(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_system_register_processes_processes FOREIGN KEY (process_id) REFERENCES processes(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE system_register_processes;

-- +goose StatementEnd
