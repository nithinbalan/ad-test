-- +goose Up
-- +goose StatementBegin

DECLARE @firstCategoryID INT;

-- Step 1: Get the first available asset category ID (smallest ID)
SELECT @firstCategoryID = MIN(id) FROM asset_categories;

-- Step 2: Ensure we have a valid category ID before proceeding
IF @firstCategoryID IS NOT NULL
BEGIN
    -- Update all related tables to use the first available category ID
    UPDATE asset_register 
    SET asset_category_id = @firstCategoryID 
    WHERE asset_category_id IS NOT NULL AND asset_category_id <> @firstCategoryID;

    UPDATE asset_criticality_matrix_asset_categories 
    SET asset_category_id = @firstCategoryID 
    WHERE asset_category_id IS NOT NULL AND asset_category_id <> @firstCategoryID;

    UPDATE asset_sub_categories 
    SET asset_category_id = @firstCategoryID 
    WHERE asset_category_id IS NOT NULL AND asset_category_id <> @firstCategoryID;

    -- Delete all asset_categories except the first available category ID (if not referenced)
    DELETE FROM asset_categories 
    WHERE id <> @firstCategoryID 
    AND id NOT IN (
        SELECT DISTINCT asset_category_id FROM asset_criticality_matrix_asset_categories
        UNION 
        SELECT DISTINCT asset_category_id FROM asset_sub_categories
    );

    -- Update the first category with new name and status
    UPDATE asset_categories 
    SET title_en = 'Information Assets', 
        title_ar = 'أصول المعلومات', 
        status = 'active' 
    WHERE id = @firstCategoryID;
END

-- Step 3: Insert remaining asset categories
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

-- Step 1: Remove newly inserted asset categories (excluding the first ID)
DELETE FROM asset_categories WHERE id > (SELECT MIN(id) FROM asset_categories);

-- Step 2: Reset the first available category ID to a generic name
UPDATE asset_categories 
SET title_en = 'General Asset', 
    title_ar = 'الأصل العام', 
    status = 'active' 
WHERE id = (SELECT MIN(id) FROM asset_categories);

-- +goose StatementEnd
