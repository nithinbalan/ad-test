-- +goose Up
-- +goose StatementBegin
CREATE TABLE document_version_clause_comments
(
    id INT PRIMARY KEY IDENTITY(1,1),
    document_version_clause_id INT NOT NULL,
    comment NVARCHAR(MAX) NOT NULL,
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT GETDATE(),
    CONSTRAINT FK_document_version_clause_comments_clause FOREIGN KEY (document_version_clause_id) REFERENCES document_version_clauses(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT FK_document_version_clause_comments_users FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE document_version_clause_comments;
-- +goose StatementEnd
