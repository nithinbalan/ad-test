-- +goose Up
-- +goose StatementBegin
-- Drop the default constraint on deadline column first
DECLARE @constraint_name NVARCHAR(200)
SELECT @constraint_name = dc.name 
FROM sys.default_constraints dc
INNER JOIN sys.columns c ON dc.parent_object_id = c.object_id AND dc.parent_column_id = c.column_id
WHERE c.name = 'is_completed' AND OBJECT_NAME(dc.parent_object_id) = 'compliance_assessment_requirements'

IF @constraint_name IS NOT NULL
BEGIN
    EXEC('ALTER TABLE compliance_assessment_requirements DROP CONSTRAINT ' + @constraint_name)
END

ALTER TABLE compliance_assessment_requirements
DROP COLUMN is_completed;

ALTER TABLE compliance_assessment_requirements
DROP COLUMN score;

ALTER TABLE compliance_assessment_requirements
DROP COLUMN comments;


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE compliance_assessment_requirements
ADD is_completed TINYINT NOT NULL DEFAULT 0;

ALTER TABLE compliance_assessment_requirements
ADD score INT NULL;

ALTER TABLE compliance_assessment_requirements
ADD comments NVARCHAR(MAX) NULL;

-- +goose StatementEnd
