-- +goose Up
-- +goose StatementBegin
CREATE TABLE sla_checks (
    id INT PRIMARY KEY IDENTITY(1, 1),
    sla_register_id INT NOT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,
    created_by INT NOT NULL,
    updated_by INT NULL,
    CONSTRAINT FK_sla_checks_sla_register FOREIGN KEY (sla_register_id) REFERENCES sla_register (id),
    CONSTRAINT FK_sla_checks_users1 FOREIGN KEY (created_by) REFERENCES users (id),
    CONSTRAINT FK_sla_checks_users2 FOREIGN KEY (updated_by) REFERENCES users (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE sla_checks;
-- +goose StatementEnd