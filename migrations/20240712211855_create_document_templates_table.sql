-- +goose Up
-- +goose StatementBegin
CREATE TABLE document_templates (
    id INT NOT NULL PRIMARY KEY IDENTITY(1,1),
    title_en NVARCHAR(100) NOT NULL,
    title_ar NVARCHAR(100) NOT NULL,
    document_type_id INT NOT NULL,
    status VARCHAR(45) NOT NULL DEFAULT 'active',
       CONSTRAINT CHK_document_templates_status CHECK (status IN ('active', 'inactive')),
    created_by_id INT,
    updated_by_id INT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME,
    CONSTRAINT FK_document_templates_document_types FOREIGN KEY (document_type_id) REFERENCES document_types(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_document_templates_users_1 FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_document_templates_users_2 FOREIGN KEY (updated_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION

);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.document_templates;
-- +goose StatementEnd
