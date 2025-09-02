-- +goose Up
-- +goose StatementBegin

-- Update the existing module from Organization Structure to Organization Unit
UPDATE modules 
SET 
    identifier = 'organization_unit',
    title_en = 'Organization Units',
    title_ar = 'الوحدة التنظيمية',
    route = '/organization/organization-units'
WHERE identifier = 'organization_structure';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Revert back to Organization Structure
UPDATE modules 
SET 
    identifier = 'organization_structure',
    title_en = 'Organization Structure',
    title_ar = 'الوحدة التنظيمية',
    route = '/organization/organization-structure'
WHERE identifier = 'organization_unit';

-- +goose StatementEnd