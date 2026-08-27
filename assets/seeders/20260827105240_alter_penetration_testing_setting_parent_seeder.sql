-- +goose Up
-- +goose StatementBegin
DECLARE @settingModuleID INT;

SELECT @settingModuleID = id
FROM modules
WHERE identifier = 'setting';

UPDATE modules
SET module_id = @settingModuleID
WHERE identifier = 'penetration_testing_setting'
  AND @settingModuleID IS NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DECLARE @penetrationTestingModuleID INT;

SELECT @penetrationTestingModuleID = id
FROM modules
WHERE identifier = 'penetration_testing';

UPDATE modules
SET module_id = @penetrationTestingModuleID
WHERE identifier = 'penetration_testing_setting'
  AND @penetrationTestingModuleID IS NOT NULL;
-- +goose StatementEnd
