-- +goose Up
-- +goose StatementBegin
ALTER TABLE bcm_execution_records 
ADD is_back_up_available_temp VARCHAR(20) NULL;
-- +goose StatementEnd

-- +goose StatementBegin
UPDATE bcm_execution_records
SET is_back_up_available_temp = CASE 
    WHEN is_back_up_available = 1 THEN 'yes' 
    ELSE 'no' 
END;
-- +goose StatementEnd

-- +goose StatementBegin
-- Drop default constraint if exists
DECLARE @ConstraintName NVARCHAR(200);
SELECT @ConstraintName = name 
FROM sys.default_constraints 
WHERE parent_object_id = OBJECT_ID('bcm_execution_records') 
AND parent_column_id = (SELECT column_id FROM sys.columns WHERE object_id = OBJECT_ID('bcm_execution_records') AND name = 'is_back_up_available');

IF @ConstraintName IS NOT NULL
    EXEC('ALTER TABLE bcm_execution_records DROP CONSTRAINT ' + @ConstraintName);

ALTER TABLE bcm_execution_records 
DROP COLUMN is_back_up_available;
-- +goose StatementEnd

-- +goose StatementBegin
EXEC sp_rename 'bcm_execution_records.is_back_up_available_temp', 'is_back_up_available', 'COLUMN';
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE bcm_execution_records 
ADD CONSTRAINT CK_bcm_execution_records_is_back_up_available 
CHECK (is_back_up_available IN ('yes', 'no', 'partial'));
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE bcm_execution_records 
DROP CONSTRAINT CK_bcm_execution_records_is_back_up_available;
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE bcm_execution_records 
ADD is_back_up_available_temp BIT DEFAULT 0;
-- +goose StatementEnd

-- +goose StatementBegin
UPDATE bcm_execution_records
SET is_back_up_available_temp = CASE 
    WHEN is_back_up_available = 'yes' THEN 1 
    ELSE 0 
END;
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE bcm_execution_records 
DROP COLUMN is_back_up_available;
-- +goose StatementEnd

-- +goose StatementBegin
EXEC sp_rename 'bcm_execution_records.is_back_up_available_temp', 'is_back_up_available', 'COLUMN';
-- +goose StatementEnd
