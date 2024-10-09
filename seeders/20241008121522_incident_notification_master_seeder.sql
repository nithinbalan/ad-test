-- +goose Up
-- +goose StatementBegin

-- Get the ID of the module with identifier 'incident'
DECLARE @module_id INT;
SET @module_id = (SELECT id FROM modules WHERE identifier = 'incident');

-- Incident module notifications
INSERT INTO notification_masters (module_id, identifier, title_en, title_ar)
VALUES
    (@module_id, 'incident', 'Notification on Incident Register activities', 'إشعار بشأن أنشطة تسجيل الحوادث'), 
    (@module_id, 'incident_action_plan', 'Notification on Incident Action Plan Activities', 'إشعار بشأن أنشطة خطة عمل الحوادث');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE notification_masters;
-- +goose StatementEnd
