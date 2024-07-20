-- +goose Up
-- +goose StatementBegin
CREATE TABLE documents (
    id INT PRIMARY KEY IDENTITY(1,1),
    reference_number VARCHAR(45) NOT NULL,
    document_type_id INT NOT NULL,
    title_en NVARCHAR(255) NOT NULL,
    title_ar NVARCHAR(255) NOT NULL,
    description_en NVARCHAR(MAX) NULL,
    description_ar NVARCHAR(MAX) NULL,
    document_category_id INT NOT NULL,
    responsible_user_id INT NOT NULL,
    responsible_unit_id INT NOT NULL,
    document_access_type_id INT NOT NULL,
    document_status_id INT NOT NULL,
    created_by_id INT NOT NULL,
    updated_by_id INT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL

    CONSTRAINT FK_documents_document_types FOREIGN KEY (document_type_id) REFERENCES document_types(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_documents_document_categories FOREIGN KEY (document_categories_id) REFERENCES document_categories(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_documents_users_1 FOREIGN KEY (responsible_user_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_documents_organization_units FOREIGN KEY (responsible_unit_id) REFERENCES organization_units(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_documents_document_access_types FOREIGN KEY (document_access_types_id) REFERENCES document_access_types(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_documents_document_statuses FOREIGN KEY (document_status_id) REFERENCES document_statuses(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_documents_users_2 FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_documents_users_3 FOREIGN KEY (updated_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.documents;
-- +goose StatementEnd
