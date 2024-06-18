-- +goose Up
-- +goose StatementBegin
CREATE TABLE risk_compliance_requirements (
    id INT PRIMARY KEY IDENTITY(1, 1),
    risk_register_id INT NOT NULL,
    compliance_requirement_id INT NOT NULL,
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL,
    CONSTRAINT FK_risk_compliance_requirements_risk_register FOREIGN KEY (risk_register_id) REFERENCES risk_register(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_risk_compliance_requirements_compliance_requirements FOREIGN KEY (compliance_requirement_id) REFERENCES compliance_requirements(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE risk_compliance_requirements;

-- +goose StatementEnd