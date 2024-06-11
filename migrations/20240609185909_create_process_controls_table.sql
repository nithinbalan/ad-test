-- +goose Up
-- +goose StatementBegin
CREATE TABLE process_controls (
    id INT PRIMARY KEY IDENTITY(1, 1),
    process_id INT NOT NULL,
    control_id INT NOT NULL,
    created_by INT NOT NULL,
    created_at DATETIME NOT NULL,
    CONSTRAINT FK_process_controls_processes FOREIGN KEY (process_id) REFERENCES processes(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_process_controls_controls FOREIGN KEY (control_id) REFERENCES controls(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE process_controls;

-- +goose StatementEnd
