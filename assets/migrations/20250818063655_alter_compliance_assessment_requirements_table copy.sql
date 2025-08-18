-- +goose Up
-- +goose StatementBegin
-- Drop the default constraint on is_completed column first
DECLARE @constraint_name NVARCHAR(200)
SELECT @constraint_name = dc.name 
FROM sys.default_constraints dc
INNER JOIN sys.columns c ON dc.parent_object_id = c.object_id AND dc.parent_column_id = c.column_id
WHERE c.name = 'is_completed' AND OBJECT_NAME(dc.parent_object_id) = 'compliance_assessment_requirements'

IF @constraint_name IS NOT NULL
BEGIN
    EXEC('ALTER TABLE compliance_assessment_requirements DROP CONSTRAINT ' + @constraint_name)
END

-- Rename the column instead of dropping
EXEC sp_rename 'compliance_assessment_requirements.is_completed', 'is_all_completed', 'COLUMN';

-- Change score column type to DECIMAL(10,2)
ALTER TABLE compliance_assessment_requirements
ALTER COLUMN score DECIMAL(10,2) NULL;

ALTER TABLE compliance_assessment_requirements
DROP COLUMN comments;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Rename column back to is_completed
EXEC sp_rename 'compliance_assessment_requirements.is_all_completed', 'is_completed', 'COLUMN';

-- Revert score column back to INT
ALTER TABLE compliance_assessment_requirements
ALTER COLUMN score INT NULL;

ALTER TABLE compliance_assessment_requirements
ADD comments NVARCHAR(MAX) NULL;

-- Restore default constraint for is_completed
ALTER TABLE compliance_assessment_requirements
ADD CONSTRAINT DF_compliance_assessment_requirements_is_completed DEFAULT 0 FOR is_completed;
-- +goose StatementEnd
