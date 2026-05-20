-- +goose Up
-- +goose StatementBegin
-- Drop existing check constraint
ALTER TABLE kri_register DROP CONSTRAINT CHK_kri_register_calculation_type;

-- Update calculation_type values from 'automatic' to 'direct' and 'manual' to 'ratio'
UPDATE kri_register SET calculation_type = 'direct' WHERE calculation_type = 'automatic';
UPDATE kri_register SET calculation_type = 'ratio' WHERE calculation_type = 'manual';

-- Add new check constraint with new values
ALTER TABLE kri_register ADD CONSTRAINT CHK_kri_register_calculation_type CHECK (calculation_type IN ('direct', 'ratio'));

-- Add new boolean fields
ALTER TABLE kri_register ADD is_lower_better BIT NOT NULL CONSTRAINT DF_kri_register_is_lower_better DEFAULT 0;
ALTER TABLE kri_register ADD is_higher_better BIT NOT NULL CONSTRAINT DF_kri_register_is_higher_better DEFAULT 0;

-- Change current_value from INT to DECIMAL(18,2)
ALTER TABLE kri_register ALTER COLUMN current_value DECIMAL(18,2) NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Remove new boolean fields
ALTER TABLE kri_register DROP CONSTRAINT DF_kri_register_is_higher_better;
ALTER TABLE kri_register DROP CONSTRAINT DF_kri_register_is_lower_better;

ALTER TABLE kri_register DROP COLUMN is_higher_better;
ALTER TABLE kri_register DROP COLUMN is_lower_better;

-- Revert current_value from DECIMAL(18,2) to INT
ALTER TABLE kri_register ALTER COLUMN current_value INT NULL;

-- Drop check constraint
ALTER TABLE kri_register DROP CONSTRAINT CHK_kri_register_calculation_type;

-- Revert calculation_type values
UPDATE kri_register SET calculation_type = 'automatic' WHERE calculation_type = 'direct';
UPDATE kri_register SET calculation_type = 'manual' WHERE calculation_type = 'ratio';

-- Add back original check constraint
ALTER TABLE kri_register ADD CONSTRAINT CHK_kri_register_calculation_type CHECK (calculation_type IN ('automatic', 'manual'));
-- +goose StatementEnd
