-- +goose Up
-- +goose StatementBegin
DECLARE @module_id INT;
SET @module_id = (SELECT id FROM modules WHERE identifier = 'incident');

-- Risk module notifications
INSERT INTO notification_masters (module_id, identifier, title_en, title_ar)
VALUES
    (@module_id, 'incident_register', 'Notification on Incident Register Activities', 'إشعار بشأن أنشطة تسجيل الحوادث');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE notification_masters;
-- +goose StatementEnd
