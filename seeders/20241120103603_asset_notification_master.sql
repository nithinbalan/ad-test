-- +goose Up
-- +goose StatementBegin

-- Get the ID of the module with identifier 'risk'
DECLARE @module_id INT;
SET @module_id = (SELECT id FROM modules WHERE identifier = 'asset');

-- Risk module notifications
INSERT INTO notification_masters (module_id, identifier, title_en, title_ar)
VALUES
    (@module_id, 'asset_register', 'Notification on Asset Register Activities', 'إشعار بشأن أنشطة سجل الأصول'),
    (@module_id, 'threat', 'Notification on Threat Activities', 'إشعار بشأن أنشطة التهديدات'),
    (@module_id, 'threat_action_plan', 'Notification on Threat Action Plan Activities', 'إشعار بشأن أنشطة خطة عمل التهديدات'),
    (@module_id, 'vulnerability', 'Notification on Vulnerability Activities', 'إشعار بشأن أنشطة نقاط الضعف'),
    (@module_id, 'vulnerability_action_plan', 'Notification on Vulnerability Action Plan Activities', 'إشعار بشأن أنشطة خطة عمل نقاط الضعف');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE notification_masters;
-- +goose StatementEnd