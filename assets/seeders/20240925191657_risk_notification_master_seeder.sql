-- +goose Up
-- +goose StatementBegin

-- Get the ID of the module with identifier 'risk'
DECLARE @module_id INT;
SET @module_id = (SELECT id FROM modules WHERE identifier = 'risk');

-- Risk module notifications
INSERT INTO notification_masters (module_id, identifier, title_en, title_ar)
VALUES
    (@module_id, 'risk', 'Notification on Risk activities', 'إشعار بشأن أنشطة المخاطر'),
    (@module_id, 'risk_action_plan', 'Notification on Risk Action Plan Activities', 'إشعار بشأن أنشطة خطة عمل المخاطر');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE notification_masters;
-- +goose StatementEnd