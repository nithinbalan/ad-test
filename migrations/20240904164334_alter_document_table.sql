-- +goose Up
-- +goose StatementBegin
ALTER TABLE documents 
ALTER COLUMN reference_number VARCHAR(45) NULL;

ALTER TABLE documents 
ADD published_by_id INT NULL,
    published_at DATETIME NULL;

ALTER TABLE documents 
ADD CONSTRAINT FK_documents_users_6 FOREIGN KEY (published_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE documents 
ALTER COLUMN reference_number VARCHAR(45) NOT NULL;

ALTER TABLE documents 
DROP CONSTRAINT FK_documents_users_6

ALTER TABLE documents 
DROP COLUMN published_by_id,
DROP COLUMN published_at
-- +goose StatementEnd
