-- +goose Up
-- +goose StatementBegin

DECLARE @settingModuleID INT;
SET @settingModuleID = (SELECT id FROM modules WHERE identifier = 'setting');


UPDATE modules
SET [order] = 1, [module_id] = @settingModuleID
WHERE identifier = 'general_setting';

UPDATE modules
SET [order] = 2
WHERE identifier = 'organization_setting';

UPDATE modules
SET [order] = 3
WHERE identifier = 'compliance_setting';

UPDATE modules
SET [order] = 4
WHERE identifier = 'risk_setting';

UPDATE modules
SET [order] = 5
WHERE identifier = 'audit_setting';

UPDATE modules
SET [order] = 6
WHERE identifier = 'document_setting';

UPDATE modules
SET [order] = 7
WHERE identifier = 'incident_setting';

UPDATE modules
SET [order] = 8
WHERE identifier = 'review_meeting_setting';

UPDATE modules
SET [order] = 9
WHERE identifier = 'audit_insight_setting';

UPDATE modules
SET [order] = 10
WHERE identifier = 'business_continuity_setting';

UPDATE modules
SET [order] = 11
WHERE identifier = 'asset_setting';

UPDATE modules
SET [order] = 13
WHERE identifier = 'common_setting';

INSERT INTO modules (identifier, title_en, title_ar, module_id, [order], icon, route)
VALUES 
    ('log', 'Logs', 'سجلات', @settingModuleID, 12, NULL, '/setting/log');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE modules;

-- +goose StatementEnd
