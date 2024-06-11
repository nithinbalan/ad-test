-- +goose Up
-- +goose StatementBegin
CREATE TABLE risk_processes (
    id INT PRIMARY KEY IDENTITY(1, 1),
    risk_register_id INT NOT NULL,
    process_id INT NOT NULL,
    created_by INT NOT NULL,
    created_at DATETIME NOT NULL,
    CONSTRAINT FK_risk_processes_risk_register FOREIGN KEY (risk_register_id) REFERENCES risk_register(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_risk_processes_processes FOREIGN KEY (process_id) REFERENCES processes(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE risk_processes;

-- +goose StatementEnd
