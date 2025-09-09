-- +goose Up
-- +goose StatementBegin
-- 1. Drop old default constraint if it exists
DECLARE @ConstraintName NVARCHAR(200);
SELECT @ConstraintName = dc.name
FROM sys.default_constraints dc
JOIN sys.columns c 
    ON dc.parent_object_id = c.object_id 
    AND dc.parent_column_id = c.column_id
WHERE dc.parent_object_id = OBJECT_ID('compliance_requirements') 
  AND c.name = 'score';
 
IF @ConstraintName IS NOT NULL
BEGIN
    EXEC('ALTER TABLE compliance_requirements DROP CONSTRAINT ' + @ConstraintName);
END
 
-- 2. Alter column type
ALTER TABLE compliance_requirements 
ALTER COLUMN score DECIMAL(10,2) NULL;

-- 3. Alter column compliance_status_id int null
ALTER TABLE compliance_requirements
ALTER COLUMN compliance_status_id INT NULL;
 
-- +goose StatementEnd
 
-- +goose Down
-- +goose StatementBegin
-- 1. Revert column type to INT
ALTER TABLE compliance_requirements 
ALTER COLUMN score INT NULL;
 
-- 2. Re-add old default constraint
ALTER TABLE compliance_requirements 
ADD CONSTRAINT DF_compliance_requirements_score DEFAULT (0) FOR score;

-- 3. Re-add old compliance_status_id NOT NULL
ALTER TABLE compliance_requirements
ALTER COLUMN compliance_status_id INT NULL;
-- +goose StatementEnd