-- +goose Up
-- +goose StatementBegin

DECLARE @module_id INT;
SET @module_id = (SELECT id FROM modules WHERE identifier = 'incident');

INSERT INTO notification_masters (module_id, identifier, title_en, title_ar)
VALUES
    (@module_id, 'incident_closure', 'Notification on Incident Closure', 'إشعار بشأن إغلاق الحادث');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM notification_masters WHERE identifier = 'incident_closure';
-- +goose StatementEnd
