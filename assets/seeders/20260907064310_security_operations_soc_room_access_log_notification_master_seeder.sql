-- +goose Up
-- +goose StatementBegin
DECLARE @moduleID INT;
SET @moduleID = (SELECT id FROM modules WHERE identifier = 'security_operations');

INSERT INTO notification_masters (module_id, identifier, title_en, title_ar)
VALUES (
    @moduleID,
    'soc_room_access_log',
    'Notification on SOC Room Access Log Activities',
    N'إشعار بشأن أنشطة سجل الدخول إلى غرفة مركز العمليات الأمنية'
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM notification_masters
WHERE identifier = 'soc_room_access_log'
  AND module_id = (SELECT id FROM modules WHERE identifier = 'security_operations');
-- +goose StatementEnd
