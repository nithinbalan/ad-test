-- +goose Up
-- +goose StatementBegin
ALTER TABLE document_versions
ALTER COLUMN url VARCHAR(1000) NULL;

ALTER TABLE document_versions
ALTER COLUMN name NVARCHAR(255) NULL;

ALTER TABLE document_versions
ALTER COLUMN size INT NULL;

ALTER TABLE document_versions
ALTER COLUMN token VARCHAR(255) NULL;

ALTER TABLE document_versions
ALTER COLUMN type VARCHAR(255) NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE document_versions
ALTER COLUMN url VARCHAR(1000) NOT NULL;

ALTER TABLE document_versions
ALTER COLUMN name NVARCHAR(255) NOT NULL;

ALTER TABLE document_versions
ALTER COLUMN size INT NOT NULL;

ALTER TABLE document_versions
ALTER COLUMN token VARCHAR(255) NOT NULL;

ALTER TABLE document_versions
ALTER COLUMN type VARCHAR(255) NOT NULL;
-- +goose StatementEnd
