-- +goose Up
-- +goose StatementBegin

-- Declare variable for constraint name
DECLARE @ConstraintName NVARCHAR(128);

-- Find the unique constraint name for reference_number
SELECT @ConstraintName = name
FROM sys.indexes
WHERE object_id = OBJECT_ID('risk_register')
AND is_unique_constraint = 1
AND EXISTS (
    SELECT 1 FROM sys.index_columns ic
    WHERE ic.object_id = indexes.object_id
    AND ic.index_id = indexes.index_id
    AND ic.column_id = (
        SELECT column_id FROM sys.columns
        WHERE object_id = OBJECT_ID('risk_register')
        AND name = 'reference_number'
    )
);

-- Drop the unique constraint if it exists
IF @ConstraintName IS NOT NULL
BEGIN
    EXEC('ALTER TABLE risk_register DROP CONSTRAINT ' + @ConstraintName);
END

-- Make column nullable
ALTER TABLE risk_register
ALTER COLUMN reference_number NVARCHAR(45) NULL;

-- Add column
ALTER TABLE risk_register
ADD risk_register_update_plan_id INT NULL;

-- Add FK
ALTER TABLE risk_register
ADD CONSTRAINT FK_risk_register_update_plan
FOREIGN KEY (risk_register_update_plan_id)
REFERENCES risk_register_update_plans(id)
ON DELETE NO ACTION
ON UPDATE NO ACTION;

-- Re-add unique constraint with explicit name
ALTER TABLE risk_register
ADD CONSTRAINT UQ_risk_register_reference_number UNIQUE (reference_number);

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

ALTER TABLE risk_register
DROP CONSTRAINT FK_risk_register_update_plan;

ALTER TABLE risk_register
DROP COLUMN risk_register_update_plan_id;

-- Drop unique constraint
ALTER TABLE risk_register
DROP CONSTRAINT UQ_risk_register_reference_number;

-- Revert column to NOT NULL
ALTER TABLE risk_register
ALTER COLUMN reference_number NVARCHAR(45) NOT NULL;

-- Re-add unique constraint
ALTER TABLE risk_register
ADD CONSTRAINT UQ_risk_register_reference_number UNIQUE (reference_number);

-- +goose StatementEnd