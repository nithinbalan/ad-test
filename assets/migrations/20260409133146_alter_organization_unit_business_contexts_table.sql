-- +goose Up
-- +goose StatementBegin

-- Drop column other_field_header
ALTER TABLE organization_unit_business_contexts
DROP COLUMN other_field_header;

-- Rename other_field_description to annual_objective_description
EXEC sp_rename 
    'organization_unit_business_contexts.other_field_description',
    'annual_objective_description',
    'COLUMN';

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

-- Rename annual_objective_description back to other_field_description
EXEC sp_rename 
    'organization_unit_business_contexts.annual_objective_description',
    'other_field_description',
    'COLUMN';

-- Add back other_field_header
ALTER TABLE organization_unit_business_contexts
ADD other_field_header NVARCHAR(MAX) NULL;

-- +goose StatementEnd