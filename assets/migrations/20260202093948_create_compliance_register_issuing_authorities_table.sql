-- +goose Up
-- +goose StatementBegin
CREATE TABLE compliance_register_issuing_authorities (
    id INT PRIMARY KEY IDENTITY(1,1),
    compliance_register_id INT NOT NULL,
    issuing_authority_id INT NOT NULL,
    created_at DATETIME NOT NULL,
    created_by_id INT NOT NULL,
    CONSTRAINT FK_compliance_register_issuing_authorities_compliance_registers FOREIGN KEY (compliance_register_id) REFERENCES compliance_register (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_compliance_register_issuing_authorities_issuing_authorities FOREIGN KEY (issuing_authority_id) REFERENCES compliance_issuing_authorities (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_compliance_register_issuing_authorities_users FOREIGN KEY (created_by_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE compliance_register_issuing_authorities;
-- +goose StatementEnd