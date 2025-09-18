-- +goose Up
-- +goose StatementBegin
ALTER TABLE document_version_clause_comments 
ADD document_version_clause_comment_id INT;

ALTER TABLE document_version_clause_comments
ADD CONSTRAINT FK_document_version_clause_comments_parent_comment
FOREIGN KEY (document_version_clause_comment_id) REFERENCES document_version_clause_comments(id)
ON DELETE NO ACTION ON UPDATE NO ACTION;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE document_version_clause_comments
DROP CONSTRAINT FK_document_version_clause_comments_parent_comment;

ALTER TABLE document_version_clause_comments
DROP COLUMN document_version_clause_comment_id;
-- +goose StatementEnd
