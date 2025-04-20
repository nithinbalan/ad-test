-- +goose Up
-- +goose StatementBegin
 
-- Drop default constraint on calculation_type (if exists)
DECLARE @calc_constraint NVARCHAR(200)
SELECT @calc_constraint = dc.name
FROM sys.default_constraints dc
JOIN sys.columns c ON c.default_object_id = dc.object_id
WHERE dc.parent_object_id = OBJECT_ID('kri_library') AND c.name = 'calculation_type';
IF @calc_constraint IS NOT NULL
    EXEC('ALTER TABLE kri_library DROP CONSTRAINT [' + @calc_constraint + ']');
 
-- Drop default constraint on status (if exists)
DECLARE @status_constraint NVARCHAR(200)
SELECT @status_constraint = dc.name
FROM sys.default_constraints dc
JOIN sys.columns c ON c.default_object_id = dc.object_id
WHERE dc.parent_object_id = OBJECT_ID('kri_library') AND c.name = 'status';
IF @status_constraint IS NOT NULL
    EXEC('ALTER TABLE kri_library DROP CONSTRAINT [' + @status_constraint + ']');
 
-- Drop CHECK constraints (if exist)
IF EXISTS (SELECT * FROM sys.check_constraints WHERE name = 'CHK_kri_library_calculation_type')
    ALTER TABLE kri_library DROP CONSTRAINT CHK_kri_library_calculation_type;
 
IF EXISTS (SELECT * FROM sys.check_constraints WHERE name = 'CHK_kri_library_status')
    ALTER TABLE kri_library DROP CONSTRAINT CHK_kri_library_status;
 
-- Drop Foreign Keys (if exist)
IF EXISTS (SELECT * FROM sys.foreign_keys WHERE name = 'FK_kri_library_kri_categories')
    ALTER TABLE kri_library DROP CONSTRAINT FK_kri_library_kri_categories;
 
IF EXISTS (SELECT * FROM sys.foreign_keys WHERE name = 'FK_kri_library_review_frequencies')
    ALTER TABLE kri_library DROP CONSTRAINT FK_kri_library_review_frequencies;
 
-- Drop old columns (if exist)
IF EXISTS (SELECT * FROM sys.columns WHERE object_id = OBJECT_ID('kri_library') AND name = 'kri_category_id')
    ALTER TABLE kri_library DROP COLUMN kri_category_id;
 
IF EXISTS (SELECT * FROM sys.columns WHERE object_id = OBJECT_ID('kri_library') AND name = 'review_frequency_id')
    ALTER TABLE kri_library DROP COLUMN review_frequency_id;
 
IF EXISTS (SELECT * FROM sys.columns WHERE object_id = OBJECT_ID('kri_library') AND name = 'calculation_type')
    ALTER TABLE kri_library DROP COLUMN calculation_type;
 
IF EXISTS (SELECT * FROM sys.columns WHERE object_id = OBJECT_ID('kri_library') AND name = 'source')
    ALTER TABLE kri_library DROP COLUMN source;
 
IF EXISTS (SELECT * FROM sys.columns WHERE object_id = OBJECT_ID('kri_library') AND name = 'status')
    ALTER TABLE kri_library DROP COLUMN status;
 
-- Add new columns
IF NOT EXISTS (SELECT * FROM sys.columns WHERE object_id = OBJECT_ID('kri_library') AND name = 'kri_category')
    ALTER TABLE kri_library ADD kri_category NVARCHAR(255) NOT NULL;
 
IF NOT EXISTS (SELECT * FROM sys.columns WHERE object_id = OBJECT_ID('kri_library') AND name = 'kri_unit')
    ALTER TABLE kri_library ADD kri_unit NVARCHAR(45) NULL;
 
-- Alter column title length
ALTER TABLE kri_library ALTER COLUMN title NVARCHAR(500) NOT NULL;
 
-- +goose StatementEnd
 
-- +goose Down
-- +goose StatementBegin
ALTER TABLE kri_library DROP COLUMN kri_category;
ALTER TABLE kri_library DROP COLUMN kri_unit;
ALTER TABLE kri_library ALTER COLUMN title NVARCHAR(255) NOT NULL;
-- +goose StatementEnd