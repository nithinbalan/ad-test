-- +goose Up
-- +goose StatementBegin
ALTER TABLE document_versions 
ADD enable_acknowledge_comment BIT NOT NULL DEFAULT 0;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Drop default constraint first
DECLARE @ConstraintName NVARCHAR(200);
SELECT @ConstraintName = dc.name
FROM sys.default_constraints dc
JOIN sys.columns c ON dc.parent_object_id = c.object_id AND dc.parent_column_id = c.column_id
WHERE OBJECT_NAME(dc.parent_object_id) = 'document_versions'
  AND c.name = 'enable_acknowledge_comment';

IF @ConstraintName IS NOT NULL
BEGIN
    EXEC('ALTER TABLE document_versions DROP CONSTRAINT ' + @ConstraintName);
END

-- Now drop the column
ALTER TABLE document_versions DROP COLUMN enable_acknowledge_comment;
-- +goose StatementEnd