-- +goose Up
-- +goose StatementBegin
CREATE TABLE project_supporting_document_documents (
    id INT PRIMARY KEY IDENTITY(1,1),
    project_supporting_document_id INT NOT NULL,
    url VARCHAR(255) NULL,
    name VARCHAR(100) NOT NULL,
    size INT NOT NULL,
    token VARCHAR(MAX) NOT NULL,
    original_name NVARCHAR(255) NULL,
    type VARCHAR(100) NOT NULL,
    created_by_id INT NULL,
    updated_by_id INT NULL,
    created_at DATETIME NOT NULL DEFAULT GETDATE(),
    updated_at DATETIME NULL,

    CONSTRAINT FK_project_supporting_document_documents_supporting_documents FOREIGN KEY (project_supporting_document_id) 
        REFERENCES project_supporting_documents(id) ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_project_supporting_document_documents_users_1 FOREIGN KEY (created_by_id) 
        REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_project_supporting_document_documents_users_2 FOREIGN KEY (updated_by_id) 
        REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE project_supporting_document_documents;
-- +goose StatementEnd
