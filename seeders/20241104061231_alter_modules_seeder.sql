-- +goose Up
-- +goose StatementBegin
-- Insert Setting Sub-Modules
DECLARE @settingModuleID INT;
SET @settingModuleID = (SELECT id FROM modules WHERE identifier = 'setting');

INSERT INTO modules (identifier, title_en, title_ar, module_id, [order], icon, route)
VALUES 
    ('asset_setting', 'Asset Setting', 'إعداد الأصول', @settingModuleID, 12, NULL, '/setting/asset-setting');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM modules
WHERE identifier = 'asset_setting';
-- +goose StatementEnd
