-- +goose Up
-- +goose StatementBegin
ALTER TABLE asset_register
ADD is_archived BIT NOT NULL DEFAULT 0;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Drop default constraint first
DECLARE @ConstraintName NVARCHAR(200);
SELECT @ConstraintName = dc.name
FROM sys.default_constraints dc
JOIN sys.columns c ON dc.parent_object_id = c.object_id AND dc.parent_column_id = c.column_id
WHERE OBJECT_NAME(dc.parent_object_id) = 'asset_register'
  AND c.name = 'is_archived';

IF @ConstraintName IS NOT NULL
BEGIN
    EXEC('ALTER TABLE asset_register DROP CONSTRAINT ' + @ConstraintName);
END

-- Now drop the column
ALTER TABLE asset_register DROP COLUMN is_archived;
-- +goose StatementEnd