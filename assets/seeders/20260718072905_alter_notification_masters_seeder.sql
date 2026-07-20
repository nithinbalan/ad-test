-- +goose Up
-- +goose StatementBegin

-- Get the ID of the module with identifier 'incident'
DECLARE @module_id INT;
SET @module_id = (SELECT id FROM modules WHERE identifier = 'organization');

-- Incident module notifications
INSERT INTO notification_masters (module_id, identifier, title_en, title_ar)
VALUES 
    (@module_id, 'project_management', 'Notification on Project Management Activities', 'إشعار بشأن أنشطة إدارة المشروع'),
    (@module_id, 'project_action_plan', 'Notification on Project Action Plan Activities', 'إشعار بشأن أنشطة خطة عمل المشروع');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM notification_masters WHERE identifier = 'project_management';
DELETE FROM notification_masters WHERE identifier = 'project_action_plan';
-- +goose StatementEnd