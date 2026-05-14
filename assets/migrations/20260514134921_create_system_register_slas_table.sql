-- +goose Up
-- +goose StatementBegin
CREATE TABLE system_register_slas (
    id INT PRIMARY KEY IDENTITY(1, 1),
    system_register_id INT NOT NULL,
    sla_register_id INT NOT NULL,
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL,
    CONSTRAINT FK_system_register_slas_system_register FOREIGN KEY (system_register_id) REFERENCES system_registers(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_system_register_slas_sla_register FOREIGN KEY (sla_register_id) REFERENCES sla_register(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE system_register_slas;

-- +goose StatementEnd
