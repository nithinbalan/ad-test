-- +goose Up
-- +goose StatementBegin
DECLARE @moduleID INT;

SET @moduleID = (SELECT id FROM modules WHERE identifier = 'project');

INSERT INTO notification_masters (module_id, identifier, title_en, title_ar)
VALUES (
    @moduleID,
    'project_security_gate',
    'Notification on Project Security Gate Activities',
    N'إشعار بشأن أنشطة بوابة أمان المشروع'
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM notification_masters
WHERE identifier = 'project_security_gate'
  AND module_id = (SELECT id FROM modules WHERE identifier = 'project');
-- +goose StatementEnd
