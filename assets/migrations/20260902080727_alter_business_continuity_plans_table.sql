-- +goose Up
-- +goose StatementBegin
-- Drop unique constraint on reference_number if it exists
DECLARE @ConstraintName NVARCHAR(200);

-- Check in sys.key_constraints for named unique constraints
SELECT @ConstraintName = kc.name
FROM sys.key_constraints kc
INNER JOIN sys.index_columns ic ON kc.parent_object_id = ic.object_id AND kc.unique_index_id = ic.index_id
INNER JOIN sys.columns c ON ic.object_id = c.object_id AND ic.column_id = c.column_id
WHERE kc.parent_object_id = OBJECT_ID('business_continuity_plans')
AND kc.type = 'UQ'
AND c.name = 'reference_number';

-- If not found, check in sys.indexes for unique indexes
IF @ConstraintName IS NULL
BEGIN
    SELECT @ConstraintName = i.name
    FROM sys.indexes i
    INNER JOIN sys.index_columns ic ON i.object_id = ic.object_id AND i.index_id = ic.index_id
    INNER JOIN sys.columns c ON ic.object_id = c.object_id AND ic.column_id = c.column_id
    WHERE i.object_id = OBJECT_ID('business_continuity_plans')
    AND i.is_unique = 1
    AND i.is_primary_key = 0
    AND c.name = 'reference_number';
END

-- Drop the constraint if found
IF @ConstraintName IS NOT NULL
BEGIN
    EXEC('ALTER TABLE business_continuity_plans DROP CONSTRAINT ' + @ConstraintName);
END
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Re-add unique constraint on reference_number
ALTER TABLE business_continuity_plans ADD CONSTRAINT UQ_business_continuity_plans_reference_number UNIQUE (reference_number);
-- +goose StatementEnd
