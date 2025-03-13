-- +goose Up
-- +goose StatementBegin
ALTER TABLE document_requirements
ADD 
    created_by_id INT NOT NULL,
    updated_by_id INT NULL,
    created_at  DATETIME NOT NULL,
    updated_at DATETIME NULL,
    CONSTRAINT FK_document_requirements_users_created FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_document_requirements_users_updated FOREIGN KEY (updated_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE document_requirements
DROP CONSTRAINT FK_document_requirements_users_created,
DROP CONSTRAINT FK_document_requirements_users_updated;

ALTER TABLE document_requirements
DROP COLUMN created_by_id,
DROP COLUMN updated_by_id,
DROP COLUMN created_at,
DROP COLUMN updated_at;
-- +goose StatementEnd
