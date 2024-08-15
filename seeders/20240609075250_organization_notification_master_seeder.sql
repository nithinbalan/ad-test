-- +goose Up
-- +goose StatementBegin

-- Get the ID of the module with identifier 'organization'
DECLARE @module_id INT;
SET @module_id = (SELECT id FROM modules WHERE identifier = 'organization');

-- Process notifications
INSERT INTO notification_masters (module_id, identifier, title_en, title_ar)
VALUES
    (@module_id, 'process', 'Notification on Process activities', 'الإشعارات لأنشطة العملية'),
    (@module_id, 'process_action_plan', 'Notification on Process Action Plan Activities', 'الإشعارات لأنشطة خطة عمل العملية');

-- Control notifications
INSERT INTO notification_masters (module_id, identifier, title_en, title_ar)
VALUES
    (@module_id, 'control', 'Notification on Control activities', 'الإشعارات لأنشطة التحكم'),
    (@module_id, 'control_action_plan', 'Notification on Control Action Plan Activities', 'الإشعارات لأنشطة خطة عمل التحكم');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE notification_masters;
-- +goose StatementEnd
