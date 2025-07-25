-- +goose Up
-- +goose StatementBegin
ALTER TABLE tprm_assessment_questions
ADD is_completed TINYINT NOT NULL DEFAULT 0;

ALTER TABLE tprm_assessment_questions
ADD completed_at DATETIME NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Drop default constraint first
DECLARE @ConstraintName NVARCHAR(200);
SELECT @ConstraintName = dc.name
FROM sys.default_constraints dc
JOIN sys.columns c ON dc.parent_object_id = c.object_id AND dc.parent_column_id = c.column_id
WHERE OBJECT_NAME(dc.parent_object_id) = 'tprm_assessment_questions'
  AND c.name = 'is_completed';

IF @ConstraintName IS NOT NULL
BEGIN
    EXEC('ALTER TABLE tprm_assessment_questions DROP CONSTRAINT ' + @ConstraintName);
END

-- Now drop the column
ALTER TABLE tprm_assessment_questions DROP COLUMN is_completed;

ALTER TABLE tprm_assessment_questions
DROP COLUMN completed_at;
-- +goose StatementEnd
