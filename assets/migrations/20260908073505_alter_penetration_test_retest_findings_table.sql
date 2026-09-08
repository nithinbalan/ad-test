-- +goose Up
-- +goose StatementBegin
DECLARE @ConstraintName NVARCHAR(200);
SELECT @ConstraintName = dc.name
FROM sys.default_constraints dc
JOIN sys.columns c ON dc.parent_object_id = c.object_id AND dc.parent_column_id = c.column_id
WHERE OBJECT_NAME(dc.parent_object_id) = 'penetration_test_retest_findings'
  AND c.name = 'is_verified';

IF @ConstraintName IS NOT NULL
BEGIN
    EXEC('ALTER TABLE penetration_test_retest_findings DROP CONSTRAINT ' + @ConstraintName);
END

ALTER TABLE penetration_test_retest_findings
DROP COLUMN is_verified;
-- +goose StatementEnd

-- +goose StatementBegin
EXEC sp_rename 'penetration_test_retest_findings.verified_at', 'result_at', 'COLUMN';
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE penetration_test_retest_findings
ADD comment NVARCHAR(MAX) NULL;
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE penetration_test_retest_findings
ADD result VARCHAR(20) NULL;
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE penetration_test_retest_findings
ADD CONSTRAINT CHK_penetration_test_retest_findings_result
CHECK (result IN ('pass', 'fail'));
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE penetration_test_retest_findings
DROP CONSTRAINT CHK_penetration_test_retest_findings_result;
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE penetration_test_retest_findings
DROP COLUMN result;
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE penetration_test_retest_findings
DROP COLUMN comment;
-- +goose StatementEnd

-- +goose StatementBegin
EXEC sp_rename 'penetration_test_retest_findings.result_at', 'verified_at', 'COLUMN';
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE penetration_test_retest_findings
ADD is_verified BIT NOT NULL CONSTRAINT DF_penetration_test_retest_findings_is_verified DEFAULT 0;
-- +goose StatementEnd
