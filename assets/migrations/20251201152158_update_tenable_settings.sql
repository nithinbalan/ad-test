-- +goose NO TRANSACTION
-- +goose Up
-- +goose StatementBegin

-- 1. Add new status column
ALTER TABLE tenable_settings
ADD status VARCHAR(10) NOT NULL DEFAULT 'active';

-- 2. Add the check constraint
ALTER TABLE tenable_settings
ADD CONSTRAINT chk_tenable_status CHECK (status IN ('active', 'inactive'));

-- 3. Migrate data from is_active → status
UPDATE tenable_settings
SET status = CASE 
                WHEN is_active = 'true' THEN 'active'
                ELSE 'inactive'
            END;

-- 4. Drop default constraint on is_active
DECLARE @constraintName NVARCHAR(200);

SELECT @constraintName = dc.name
FROM sys.default_constraints dc
JOIN sys.columns c ON c.default_object_id = dc.object_id
WHERE c.object_id = OBJECT_ID('tenable_settings')
  AND c.name = 'is_active';

IF @constraintName IS NOT NULL
    EXEC('ALTER TABLE tenable_settings DROP CONSTRAINT ' + @constraintName);

-- 5. Drop the old column
ALTER TABLE tenable_settings
DROP COLUMN is_active;

-- +goose StatementEnd
