-- +goose Up
-- +goose StatementBegin

DECLARE @itAssetID INT;
DECLARE @physicalAssetID INT;
DECLARE @hrAssetID INT;

SELECT @itAssetID = id FROM asset_types WHERE title_en = 'IT Asset';

SELECT @physicalAssetID = id FROM asset_types WHERE title_en = 'Physical Asset';

SELECT @hrAssetID = id FROM asset_types WHERE title_en = 'Human Resource';

INSERT INTO asset_categories (asset_type_id, title_en, title_ar, status)
VALUES
    (@itAssetID, 'Hardware', 'الأجهزة', 'active'),
    (@itAssetID, 'Software', 'البرمجيات', 'active'),
    (@itAssetID, 'Networking Equipment', 'معدات الشبكات', 'active'),
    (@itAssetID, 'Cloud Resources', 'الموارد السحابية', 'active');

INSERT INTO asset_categories (asset_type_id, title_en, title_ar, status)
VALUES
    (@physicalAssetID, 'Buildings', 'المباني', 'active'),
    (@physicalAssetID, 'Vehicles', 'المركبات', 'active'),
    (@physicalAssetID, 'Machinery', 'الآلات', 'active'),
    (@physicalAssetID, 'Furniture', 'الأثاث', 'active');

INSERT INTO asset_categories (asset_type_id, title_en, title_ar, status)
VALUES
    (@hrAssetID, 'Employees', 'الموظفين', 'active'),
    (@hrAssetID, 'Contractors', 'المقاولون', 'active'),
    (@hrAssetID, 'Training Programs', 'برامج التدريب', 'active'),
    (@hrAssetID, 'Recruitment Resources', 'موارد التوظيف', 'active');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE asset_categories;

-- +goose StatementEnd
