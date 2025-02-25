-- +goose Up
-- +goose StatementBegin

-- Step 1: Update all references in related tables to asset_category_id = 1
UPDATE asset_register SET asset_category_id = 1 WHERE asset_category_id IS NOT NULL;
UPDATE asset_criticality_matrix_asset_categories SET asset_category_id = 1 WHERE asset_category_id IS NOT NULL;
UPDATE asset_sub_categories SET asset_category_id = 1 WHERE asset_category_id IS NOT NULL;

-- Step 2: Delete all asset_categories except id=1 (only if they are not referenced)
DELETE FROM asset_categories 
WHERE id <> 1 
AND id NOT IN (
    SELECT DISTINCT asset_category_id FROM asset_criticality_matrix_asset_categories
    UNION 
    SELECT DISTINCT asset_category_id FROM asset_sub_categories
);

-- Step 3: Update id=1 with new data
UPDATE asset_categories 
SET title_en = 'Information Assets', 
    title_ar = 'أصول المعلومات', 
    status = 'active' 
WHERE id = 1;

-- Step 4: Insert remaining asset categories
INSERT INTO asset_categories (title_en, title_ar, status) VALUES
('Data Assets', 'أصول البيانات', 'active'),
('Business Applications', 'تطبيقات الأعمال', 'active'),
('Hardware Assets', 'أصول الأجهزة', 'active'),
('Software Assets', 'أصول البرمجيات', 'active'),
('Cloud Assets', 'أصول السحابة', 'active'),
('Network Assets', 'أصول الشبكة', 'active'),
('People Assets', 'أصول الأشخاص', 'active'),
('Physical Assets', 'الأصول المادية', 'active'),
('Communication Assets', 'أصول الاتصالات', 'active');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Step 1: Remove newly inserted asset categories (except id=1)
DELETE FROM asset_categories WHERE id > 1;

-- Step 2: Reset id=1 back to a generic category
UPDATE asset_categories 
SET title_en = 'General Asset', 
    title_ar = 'الأصل العام', 
    status = 'active' 
WHERE id = 1;

-- +goose StatementEnd
