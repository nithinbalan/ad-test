
-- +goose Up
-- +goose StatementBegin
ALTER TABLE document_version_clause_comments
ADD level INT NULL;

ALTER TABLE document_version_clause_comments
ADD updated_at DATETIME NULL;

ALTER TABLE document_version_clause_comments
ADD comment_status_id INT NULL;

ALTER TABLE document_version_clause_comments
ADD CONSTRAINT FK_document_version_clause_comments_comment_statuses
FOREIGN KEY (comment_status_id) REFERENCES comment_statuses(id)
ON DELETE NO ACTION ON UPDATE NO ACTION;

ALTER TABLE document_version_clause_comments
ADD updated_by_id INT NULL;

ALTER TABLE document_version_clause_comments
ADD CONSTRAINT FK_document_version_clause_comments_user_1
FOREIGN KEY (updated_by_id) REFERENCES users(id)
ON DELETE NO ACTION ON UPDATE NO ACTION;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE document_version_clause_comments
DROP CONSTRAINT FK_document_version_clause_comments_comment_statuses;


ALTER TABLE document_version_clause_comments
DROP COLUMN comment_status_id;

ALTER TABLE document_version_clause_comments
DROP CONSTRAINT FK_document_version_clause_comments_user_1;


ALTER TABLE document_version_clause_comments
DROP COLUMN updated_by_id;

ALTER TABLE document_version_clause_comments
DROP COLUMN level;

ALTER TABLE document_version_clause_comments
DROP COLUMN updated_at;
-- +goose StatementEnd
