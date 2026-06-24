-- +goose Up
-- +goose StatementBegin
CREATE TABLE system_register_impact_types (
    id INT IDENTITY(1,1) PRIMARY KEY,
    system_register_id INT NOT NULL,
    impact_type_id INT NOT NULL,
    created_by_id INT,
    created_at DATETIME NOT NULL,
    CONSTRAINT FK_system_register_impact_types_system_registers FOREIGN KEY (system_register_id) REFERENCES system_registers(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_system_register_impact_types_impact_categories FOREIGN KEY (impact_type_id) REFERENCES impact_categories(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_system_register_impact_types_users_1 FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE system_register_impact_types;
-- +goose StatementEnd
