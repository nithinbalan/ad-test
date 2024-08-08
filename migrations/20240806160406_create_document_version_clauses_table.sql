-- +goose Up
-- +goose StatementBegin
CREATE TABLE document_version_clauses (
    id INT PRIMARY KEY IDENTITY(1,1),
    document_version_clause_id INT NULL,
    clause_number VARCHAR(45) NULL,
    clause_title_en NVARCHAR(255) NOT NULL,
    clause_title_ar NVARCHAR(255) NOT NULL,
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL,
    updated_by_id INT NULL,
    updated_at DATETIME NULL,
    CONSTRAINT FK_document_version_clause_self FOREIGN KEY (document_version_clause_id) REFERENCES document_version_clauses(id),
    CONSTRAINT FK_document_version_clauses_users_1 FOREIGN KEY (created_by_id) REFERENCES users(id),
    CONSTRAINT FK_document_version_clauses_users_2 FOREIGN KEY (updated_by_id) REFERENCES users(id)  
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE document_version_clauses;
-- +goose StatementEnd
