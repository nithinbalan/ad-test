-- +goose Up
-- +goose StatementBegin
ALTER TABLE business_continuity_plans ADD parent_id INT NULL;

-- Drop unique constraint on reference_number if it exists
IF EXISTS (
    SELECT 1 
    FROM sys.indexes 
    WHERE object_id = OBJECT_ID('business_continuity_plans') 
    AND name LIKE 'UQ%reference_number%'
)
BEGIN
    DECLARE @ConstraintName NVARCHAR(200);
    SELECT @ConstraintName = name 
    FROM sys.indexes 
    WHERE object_id = OBJECT_ID('business_continuity_plans') 
    AND name LIKE 'UQ%reference_number%';
    
    EXEC('ALTER TABLE business_continuity_plans DROP CONSTRAINT ' + @ConstraintName);
END

ALTER TABLE business_continuity_plans
ADD CONSTRAINT FK_business_continuity_plans_parent
FOREIGN KEY (parent_id) REFERENCES business_continuity_plans(id)
ON DELETE NO ACTION ON UPDATE NO ACTION;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE business_continuity_plans DROP CONSTRAINT FK_business_continuity_plans_parent;
ALTER TABLE business_continuity_plans DROP COLUMN parent_id;
ALTER TABLE business_continuity_plans ADD CONSTRAINT UQ_business_continuity_plans_reference_number UNIQUE (reference_number);
-- +goose StatementEnd
