-- +goose Up
-- +goose StatementBegin

-- Update route for GRC Dashboard sub-module
UPDATE modules
SET route = 'dashboards/grc'
WHERE identifier = 'grc_dashboard';

-- Update icon for Audit Menu
UPDATE modules
SET icon = 'icon-audit-menu.svg'
WHERE identifier = 'audit';

-- Update icon for Audit Insight Menu
UPDATE modules
SET icon = 'icon-audit-insight-menu.svg'
WHERE identifier = 'audit_insights';

-- Update icon for Incident Menu
UPDATE modules
SET icon = 'icon-incident-menu.svg'
WHERE identifier = 'incident';

-- Insert General Setting Sub-Module under Settings
INSERT INTO modules (identifier, title_en, title_ar, module_id, [order], icon, route)
VALUES 
    ('general_setting', 'General Setting', 'الإعداد العام', 
     (SELECT module_id FROM modules WHERE identifier = 'setting'), 10, NULL, '/setting/general-setting');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Revert route change for GRC Dashboard sub-module
UPDATE modules
SET route = '/dashboard/my-dashboard'
WHERE identifier = 'grc_dashboard';

-- Revert icon for Audit Menu
UPDATE modules
SET icon = 'icon-risk.svg'
WHERE identifier = 'audit';

-- Revert icon for Audit Insight Menu
UPDATE modules
SET icon = 'icon-risk.svg'
WHERE identifier = 'audit_insights';

-- Revert icon for Incident Menu
UPDATE modules
SET icon = 'icon-incident.svg'
WHERE identifier = 'incident';

-- Remove the General Setting sub-module from Settings
DELETE FROM modules
WHERE identifier = 'general_setting';

-- +goose StatementEnd