-- +goose Up
-- +goose StatementBegin
CREATE TABLE compliance_requirement_associated_documents (
    id INT PRIMARY KEY IDENTITY(1, 1),
    compliance_register_id INT NOT NULL,
    compliance_requirement_id INT NOT NULL,
    document_id INT NOT NULL,
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL,
    CONSTRAINT FK_compliance_requirement_associated_documents_compliance_register FOREIGN KEY (compliance_register_id) REFERENCES compliance_register(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_compliance_requirement_associated_documents_compliance_requirements FOREIGN KEY (compliance_requirement_id) REFERENCES compliance_requirements(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_compliance_requirement_associated_documents_documents FOREIGN KEY (document_id) REFERENCES documents(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE compliance_requirement_associated_documents;
-- +goose StatementEnd