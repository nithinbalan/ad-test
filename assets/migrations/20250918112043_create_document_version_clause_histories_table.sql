-- +goose Up
-- +goose StatementBegin
CREATE TABLE document_version_clause_histories (
    id INT PRIMARY KEY IDENTITY(1, 1),
    document_version_clause_id INT NOT NULL,
    old_value NVARCHAR(MAX) NULL,
    new_value NVARCHAR(MAX) NULL,
    changed_field NVARCHAR(MAX) NULL,
    created_at DATETIME NOT NULL,
    created_by_id INT NOT NULL,
    CONSTRAINT FK_document_version_clause_histories_clauses FOREIGN KEY (document_version_clause_id) REFERENCES document_version_clauses (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_document_version_clause_histories_users_1 FOREIGN KEY (created_by_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.document_version_clause_histories;
-- +goose StatementEnd