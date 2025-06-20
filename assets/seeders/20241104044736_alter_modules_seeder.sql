-- +goose Up
-- +goose StatementBegin
DECLARE @assetModuleID INT;

UPDATE modules
SET [order] = 12
WHERE identifier = 'setting';


-- Insert Asset Module
INSERT INTO modules (identifier, title_en, title_ar, [order], icon, route)
VALUES ('asset', 'Asset', 'أصل', 11, 'icon-asset.svg', '/asset');
SELECT @assetModuleID = SCOPE_IDENTITY();

-- Insert Asset Sub-Modules
INSERT INTO modules (identifier, title_en, title_ar, module_id, [order], icon, route)
VALUES 
    ('asset_register', 'Asset Register', 'سجل الأصول', @assetModuleID, 1, NULL, '/asset/asset-register'),
    ('asset_criticality_matrix', 'Asset Criticality Matrix', 'مصفوفة أهمية الأصول', @assetModuleID, 2, NULL, '/asset/asset-criticality-matrix'),
    ('vulnerability', 'Vulnerabilities', 'نقاط الضعف', @assetModuleID, 3, NULL, '/asset/vulnerabilities'),
    ('threat', 'Threats', 'التهديدات', @assetModuleID, 4, NULL, '/asset/threats'),
    ('risk', 'Risks', 'المخاطر', @assetModuleID, 5, NULL, '/asset/risks');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DECLARE @assetModuleID INT;

SELECT @assetModuleID = id FROM modules WHERE identifier = 'asset';

UPDATE modules 
SET [order] = 11 
WHERE identifier = 'setting';

DELETE FROM modules
WHERE module_id = @assetModuleID;

DELETE FROM modules
WHERE identifier = 'asset';
-- +goose StatementEnd
