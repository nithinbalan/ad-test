-- +goose Up
-- +goose StatementBegin
DECLARE @projectModuleID INT;

SET @projectModuleID = (SELECT id FROM modules WHERE identifier = 'project');

IF EXISTS (SELECT 1 FROM notification_masters WHERE identifier IN ('project_management', 'project_action_plan'))
BEGIN
    UPDATE notification_masters
    SET module_id = @projectModuleID
    WHERE identifier IN ('project_management', 'project_action_plan');
END
ELSE
BEGIN
    INSERT INTO notification_masters (module_id, identifier, title_en, title_ar)
    VALUES 
        (@projectModuleID, 'project_management', 'Notification on Project Management Activities', N'إشعار بشأن أنشطة إدارة المشروع'),
        (@projectModuleID, 'project_action_plan', 'Notification on Project Action Plan Activities', N'إشعار بشأن أنشطة خطة عمل المشروع');
END
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DECLARE @orgModuleID INT;

SET @orgModuleID = (SELECT id FROM modules WHERE identifier = 'organization');

UPDATE notification_masters
SET module_id = @orgModuleID
WHERE identifier IN ('project_management', 'project_action_plan');
-- +goose StatementEnd

