-- +goose Up
-- +goose StatementBegin

DECLARE @BCM INT;
DECLARE @BCMPolicy INT;
DECLARE @settingModuleID INT;

SELECT @BCM = id FROM modules WHERE identifier = 'business_continuity';
SELECT @BCMPolicy = id FROM modules WHERE identifier = 'policies' and module_id=@BCM;
SELECT @settingModuleID = id FROM modules WHERE identifier = 'setting';

-- Update identifier for BCM Policies
UPDATE modules
SET identifier = 'bcm_policies'
WHERE id = @BCMPolicy;

-- Insert Master Setting Sub-Module under Settings
INSERT INTO modules (identifier, title_en, title_ar, module_id, [order], icon, route)
VALUES 
    ('common_setting', 'Masters', 'الماجستير', @settingModuleID, 1, NULL, '/setting/master-setting');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Revert identifier for BCM Policies
UPDATE modules
SET identifier = 'policies'
WHERE module_id = @BCMPolicy;

-- Remove the General Setting sub-module from Settings
DELETE FROM modules
WHERE identifier = 'common_setting';
-- +goose StatementEnd
