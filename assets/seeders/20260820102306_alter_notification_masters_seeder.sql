-- +goose Up
-- +goose StatementBegin
DECLARE @moduleID INT;
SET @moduleID = (SELECT id FROM modules WHERE identifier = 'key_performance_indicator');

INSERT INTO notification_masters (module_id, identifier, title_en, title_ar)
VALUES (
    @moduleID,
    'kpi_register',
    'Notification on KPI Register Activities',
    N'إشعار بشأن أنشطة سجل مؤشرات الأداء الرئيسية'
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM notification_masters
WHERE identifier = 'kpi_register'
  AND module_id = (SELECT id FROM modules WHERE identifier = 'key_performance_indicator');
-- +goose StatementEnd
