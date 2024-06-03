-- +goose Up
-- +goose StatementBegin

-- Get the ID of the module with identifier 'meeting'
DECLARE @module_id INT;
SET @module_id = (SELECT id FROM modules WHERE identifier = 'meeting');

INSERT INTO notification_masters (module_id, identifier, title_en, title_ar)
VALUES
    (@module_id, 'create_meeting_action_plan', 'Notification on creation of a new Action Plan', 'إشعار بشأن إنشاء خطة عمل جديدة');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE notification_masters;
-- +goose StatementEnd
