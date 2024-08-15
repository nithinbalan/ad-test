-- +goose Up
-- +goose StatementBegin
CREATE TABLE document_template_documents (
    id INT PRIMARY KEY IDENTITY(1,1),
    document_template_id INT NOT NULL,
    url VARCHAR(1000) NOT NULL,
    name NVARCHAR(255) NOT NULL,
    size INT NOT NULL,
    token VARCHAR(255) NOT NULL,
    original_name NVARCHAR(255) NULL,
    type VARCHAR(255) NOT NULL,
    created_by_id INT,
    updated_by_id INT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,
    CONSTRAINT FK_document_template_documents_document_templates FOREIGN KEY (document_template_id) REFERENCES document_templates(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_document_template_documents_users_1 FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_document_template_documents_users_2 FOREIGN KEY (updated_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.document_template_documents;
-- +goose StatementEnd
