-- +goose Up
-- +goose StatementBegin
ALTER TABLE documents 
ADD is_internal TINYINT NOT NULL DEFAULT 0;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DECLARE @ConstraintName NVARCHAR(200);
SELECT @ConstraintName = dc.name
FROM sys.default_constraints dc
JOIN sys.columns c ON dc.parent_object_id = c.object_id AND dc.parent_column_id = c.column_id
WHERE OBJECT_NAME(dc.parent_object_id) = 'documents' AND c.name = 'is_internal';
IF @ConstraintName IS NOT NULL
BEGIN
    EXEC('ALTER TABLE documents DROP CONSTRAINT ' + @ConstraintName);
END
-- Now drop the column
ALTER TABLE documents DROP COLUMN is_internal;
-- +goose StatementEnd
