-- +goose Up
-- +goose StatementBegin
ALTER TABLE document_version_clause_comments
ALTER COLUMN document_version_clause_id INT NULL;

ALTER TABLE document_version_clause_comments 
ADD document_version_id INT;

ALTER TABLE document_version_clause_comments
ADD CONSTRAINT FK_document_version_clause_comments_document_version
FOREIGN KEY (document_version_id) REFERENCES documents(id)
ON DELETE NO ACTION ON UPDATE NO ACTION;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE document_version_clause_comments
ALTER COLUMN document_version_clause_id INT NOT NULL;

ALTER TABLE document_version_clause_comments
DROP CONSTRAINT FK_document_version_clause_comments_document_version;

ALTER TABLE document_version_clause_comments
DROP COLUMN document_version_id;
-- +goose StatementEnd
