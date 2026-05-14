-- +goose Up
-- +goose StatementBegin
CREATE TABLE system_register_recovery_locations (
    id INT PRIMARY KEY IDENTITY(1,1),
    system_register_id INT NOT NULL,
    recovery_location_id INT NOT NULL,
    created_by_id INT,
    created_at DATETIME NOT NULL,
    CONSTRAINT FK_system_register_recovery_locations_system_registers FOREIGN KEY (system_register_id) REFERENCES system_registers(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_system_register_recovery_locations_recovery_locations FOREIGN KEY (recovery_location_id) REFERENCES locations(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_system_register_recovery_locations_users_1 FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE system_register_recovery_locations;
-- +goose StatementEnd