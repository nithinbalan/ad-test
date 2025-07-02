-- +goose Up
-- +goose StatementBegin
ALTER TABLE document_acknowledge_users
ADD document_version_id INT NULL;

ALTER TABLE document_acknowledge_users
ADD CONSTRAINT FK_document_acknowledge_users_document_versions FOREIGN KEY (document_version_id) REFERENCES document_versions(id) ON DELETE NO ACTION ON UPDATE NO ACTION;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE document_acknowledge_users
DROP CONSTRAINT FK_document_acknowledge_users_document_versions;

ALTER TABLE document_acknowledge_users
DROP COLUMN document_version_id;
-- +goose StatementEnd
