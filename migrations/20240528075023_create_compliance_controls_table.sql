-- +goose Up
-- +goose StatementBegin
CREATE TABLE compliance_controls (
    id INT PRIMARY KEY IDENTITY(1, 1),
    compliance_register_id INT NOT NULL,
    compliance_requirement_id INT NULL,
    control_id INT NOT NULL,
    created_by INT NOT NULL,
    created_at DATETIME NOT NULL,
    CONSTRAINT FK_compliance_controls_compliance_register FOREIGN KEY (compliance_register_id) REFERENCES compliance_register(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_compliance_controls_compliance_requirement FOREIGN KEY (compliance_requirement_id) REFERENCES compliance_requirements(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_compliance_controls_controls FOREIGN KEY (control_id) REFERENCES controls(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE compliance_controls;

-- +goose StatementEnd
