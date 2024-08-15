-- +goose Up
-- +goose StatementBegin
CREATE TABLE document_shared_units (
    id INT PRIMARY KEY IDENTITY(1,1),
    document_id INT NOT NULL,
    organization_unit_id INT NOT NULL,
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL,

    CONSTRAINT FK_document_shared_units_documents FOREIGN KEY (document_id) REFERENCES documents(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_document_shared_units_organization_units FOREIGN KEY (organization_unit_id) REFERENCES organization_units(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_document_shared_units_users FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.document_shared_units;
-- +goose StatementEnd
