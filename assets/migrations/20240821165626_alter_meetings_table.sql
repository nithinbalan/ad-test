-- +goose Up
-- +goose StatementBegin

-- Step 1: Add the reference_number column without constraints, allowing NULLs temporarily
ALTER TABLE meetings
ADD reference_number VARCHAR(45) NULL;

-- Step 5: Add the description column
ALTER TABLE meetings
ADD description NVARCHAR(MAX) NULL;

-- Step 6: Change the data types of start_time and end_time columns
ALTER TABLE meetings
ALTER COLUMN start_time NVARCHAR(45) NULL;

ALTER TABLE meetings
ALTER COLUMN end_time NVARCHAR(45) NULL;

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin

-- Drop the reference_number and description columns
ALTER TABLE meetings
DROP COLUMN reference_number;

ALTER TABLE meetings
DROP COLUMN description;

-- Revert the data types of start_time and end_time columns
ALTER TABLE meetings
ALTER COLUMN start_time TIME NOT NULL;

ALTER TABLE meetings
ALTER COLUMN end_time TIME NOT NULL;

-- +goose StatementEnd
