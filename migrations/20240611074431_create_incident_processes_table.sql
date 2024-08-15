-- +goose Up
-- +goose StatementBegin
CREATE TABLE incident_processes (
    id INT PRIMARY KEY IDENTITY(1, 1),
    incident_register_id INT NOT NULL,
    process_id INT NOT NULL,
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL,
    CONSTRAINT FK_incident_processes_incident_register FOREIGN KEY (incident_register_id) REFERENCES incident_register(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_incident_processes_processes FOREIGN KEY (process_id) REFERENCES processes(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE incident_processes;

-- +goose StatementEnd
