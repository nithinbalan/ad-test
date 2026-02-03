-- +goose Up
-- +goose StatementBegin
CREATE TABLE compliance_register_issuing_authority_entities (
    id INT PRIMARY KEY IDENTITY(1,1),
    compliance_register_issuing_authority_id INT NOT NULL,
    circular_no NVARCHAR(500) NOT NULL,
    date_of_receiving DATETIME NOT NULL,
    hijri_date_of_receiving DATETIME NOT NULL,
    created_at DATETIME NOT NULL,
    created_by_id INT NOT NULL,
    updated_by_id INT NULL,
    updated_at DATETIME NULL,
    CONSTRAINT FK_compliance_register_issuing_authority_entities_compliance_register_issuing_authorities FOREIGN KEY (compliance_register_issuing_authority_id) REFERENCES compliance_register_issuing_authorities (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_compliance_register_issuing_authority_entities_users_1 FOREIGN KEY (created_by_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_compliance_register_issuing_authority_entities_users_2 FOREIGN KEY (updated_by_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE compliance_register_issuing_authority_entities;
-- +goose StatementEnd