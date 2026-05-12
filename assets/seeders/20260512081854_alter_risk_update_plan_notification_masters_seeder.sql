-- +goose Up
-- +goose StatementBegin
DECLARE @module_id INT;
SET @module_id = (SELECT id FROM modules WHERE identifier = 'risk');

-- Risk module notifications
INSERT INTO notification_masters (module_id, identifier, title_en, title_ar)
VALUES
    (@module_id, 'risk_register_update_plan', 'Notification on Risk Register Update Plan', 'الإخطار على تحديث خطة سجل المخاطر');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE notification_masters;
-- +goose StatementEnd

