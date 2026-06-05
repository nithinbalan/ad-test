-- +goose Up
-- +goose StatementBegin

-- Drop default constraint first
DECLARE @ConstraintName NVARCHAR(200);

SELECT @ConstraintName = dc.name
FROM sys.default_constraints dc
JOIN sys.columns c
    ON dc.parent_object_id = c.object_id
   AND dc.parent_column_id = c.column_id
WHERE OBJECT_NAME(dc.parent_object_id) = 'users'
  AND c.name = 'is_head';

IF @ConstraintName IS NOT NULL
BEGIN
    EXEC('ALTER TABLE users DROP CONSTRAINT ' + @ConstraintName);
END

ALTER TABLE users
DROP COLUMN is_head;

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

ALTER TABLE users
ADD is_head BIT NOT NULL DEFAULT 0;

-- +goose StatementEnd