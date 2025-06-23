-- +goose Up
-- +goose StatementBegin
-- Get the ID of the module with identifier 'audit'
DECLARE @module_id INT;
SET @module_id = (SELECT id FROM modules WHERE identifier = 'audit');

-- Risk module notifications
INSERT INTO notification_masters (module_id, identifier, title_en, title_ar)
VALUES
    (@module_id, 'audit', 'Notification on Audit activities', 'إشعار بشأن أنشطة التدقيق'),
    (@module_id, 'audit_finding_action_plan', 'Notification on Audit Finding Action Plan Activities', 'إشعار بشأن أنشطة خطة عمل نتائج التدقيق');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE notification_masters;
-- +goose StatementEnd
