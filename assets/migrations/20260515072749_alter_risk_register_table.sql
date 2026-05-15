-- +goose Up
-- +goose StatementBegin

-- Drop the unique constraint on reference_number to allow multiple NULLs
IF EXISTS (
    SELECT 1 FROM sys.indexes 
    WHERE name = 'UQ_risk_register_reference_number' 
    AND object_id = OBJECT_ID('risk_register')
)
BEGIN
    ALTER TABLE risk_register DROP CONSTRAINT UQ_risk_register_reference_number;
END

-- Create a filtered unique index that only enforces uniqueness on non-NULL values
CREATE UNIQUE NONCLUSTERED INDEX UQ_risk_register_reference_number_filtered
ON risk_register(reference_number)
WHERE reference_number IS NOT NULL;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Drop the filtered unique index
IF EXISTS (
    SELECT 1 FROM sys.indexes 
    WHERE name = 'UQ_risk_register_reference_number_filtered' 
    AND object_id = OBJECT_ID('risk_register')
)
BEGIN
    DROP INDEX UQ_risk_register_reference_number_filtered ON risk_register;
END

-- Re-add the original unique constraint
ALTER TABLE risk_register
ADD CONSTRAINT UQ_risk_register_reference_number UNIQUE (reference_number);

-- +goose StatementEnd
