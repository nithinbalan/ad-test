-- +goose Up
-- +goose StatementBegin

-- Get the ID of the module with identifier 'meeting'
DECLARE @module_id INT;
SET @module_id = (SELECT id FROM modules WHERE identifier = 'meeting');

-- Compliance Register notifications
INSERT INTO notification_masters (module_id, identifier, title_en, title_ar)
VALUES
    (@module_id, 'meeting', 'Notification on Meeting activities', 'إشعار بشأن أنشطة الاجتماع'),
    (@module_id, 'meeting_action_plan', 'Notification on Meeting Action Plan Activities', 'إشعار بشأن أنشطة خطة عمل الاجتماع')

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE notification_masters;
-- +goose StatementEnd
