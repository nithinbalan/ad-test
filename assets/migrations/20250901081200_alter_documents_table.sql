-- +goose Up
-- +goose StatementBegin
ALTER TABLE documents 
ADD is_folder TINYINT NOT NULL DEFAULT 0;

ALTER TABLE documents 
ADD parent_document_id INT;

ALTER TABLE documents
ADD CONSTRAINT FK_documents_parent_document
FOREIGN KEY (parent_document_id) REFERENCES documents(id)
ON DELETE NO ACTION ON UPDATE NO ACTION;

ALTER TABLE documents
ALTER COLUMN document_type_id INT NULL;

ALTER TABLE documents
ALTER COLUMN document_category_id INT NULL;

ALTER TABLE documents
ALTER COLUMN responsible_user_id INT NULL;

ALTER TABLE documents
ALTER COLUMN responsible_unit_id INT NULL;

ALTER TABLE documents
ALTER COLUMN document_status_id INT NULL;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DECLARE @ConstraintName NVARCHAR(200);
SELECT @ConstraintName = dc.name
FROM sys.default_constraints dc
JOIN sys.columns c ON dc.parent_object_id = c.object_id AND dc.parent_column_id = c.column_id
WHERE OBJECT_NAME(dc.parent_object_id) = 'documents' AND c.name = 'is_folder';
IF @ConstraintName IS NOT NULL
BEGIN
    EXEC('ALTER TABLE documents DROP CONSTRAINT ' + @ConstraintName);
END
-- Now drop the column
ALTER TABLE documents DROP COLUMN is_folder;

ALTER TABLE documents
DROP CONSTRAINT FK_documents_parent_document;

ALTER TABLE documents
DROP COLUMN parent_document_id;

ALTER TABLE documents
ALTER COLUMN document_type_id INT NOT NULL;

ALTER TABLE documents
ALTER COLUMN document_category_id INT NOT NULL;

ALTER TABLE documents
ALTER COLUMN responsible_user_id INT NOT NULL;

ALTER TABLE documents
ALTER COLUMN responsible_unit_id INT NOT NULL;

ALTER TABLE documents
ALTER COLUMN document_status_id INT NOT NULL;
-- +goose StatementEnd
