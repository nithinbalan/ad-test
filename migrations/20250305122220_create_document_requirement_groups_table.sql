-- +goose Up
-- +goose StatementBegin
CREATE TABLE document_requirement_groups (
    id INT PRIMARY KEY IDENTITY(1, 1),
    document_id INT,
    title_en NVARCHAR(250) NOT NULL,
    title_ar NVARCHAR(250) NOT NULL,    
    status VARCHAR(10) NOT NULL DEFAULT 'active',
        CONSTRAINT CHK_document_requirement_groups_status CHECK (status IN ('active', 'inactive')),
    CONSTRAINT FK_document_requirement_groups_documents FOREIGN KEY (document_id) REFERENCES documents(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.document_requirement_groups;
-- +goose StatementEnd
