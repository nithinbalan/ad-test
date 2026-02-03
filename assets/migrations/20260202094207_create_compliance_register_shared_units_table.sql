-- +goose Up
-- +goose StatementBegin
CREATE TABLE compliance_register_shared_units (
    id INT PRIMARY KEY IDENTITY(1,1),
    compliance_register_id INT NOT NULL,
    department_id INT NOT NULL,
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL,

    CONSTRAINT FK_compliance_register_shared_units_compliance_registers FOREIGN KEY (compliance_register_id) REFERENCES compliance_register(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_compliance_register_shared_units_departments FOREIGN KEY (department_id) REFERENCES organization_units(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_compliance_register_shared_units_users FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE compliance_register_shared_units;
-- +goose StatementEnd