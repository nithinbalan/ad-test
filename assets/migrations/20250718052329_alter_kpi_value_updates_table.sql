-- +goose Up
-- +goose StatementBegin
ALTER TABLE kpi_value_updates
ADD is_latest BIT NOT NULL DEFAULT 0;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Drop default constraint first
DECLARE @ConstraintName NVARCHAR(200);
SELECT @ConstraintName = dc.name
FROM sys.default_constraints dc
JOIN sys.columns c ON dc.parent_object_id = c.object_id AND dc.parent_column_id = c.column_id
WHERE OBJECT_NAME(dc.parent_object_id) = 'kpi_value_updates'
  AND c.name = 'is_latest';

IF @ConstraintName IS NOT NULL
BEGIN
    EXEC('ALTER TABLE kpi_value_updates DROP CONSTRAINT ' + @ConstraintName);
END

-- Now drop the column
ALTER TABLE kpi_value_updates DROP COLUMN is_latest;
-- +goose StatementEnd


