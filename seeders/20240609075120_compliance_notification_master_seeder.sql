-- +goose Up
-- +goose StatementBegin

-- Get the ID of the module with identifier 'compliance'
DECLARE @module_id INT;
SET @module_id = (SELECT id FROM modules WHERE identifier = 'compliance');

-- Compliance Register notifications
INSERT INTO notification_masters (module_id, identifier, title_en, title_ar)
VALUES
    (@module_id, 'compliance_register', 'Enable/Disable notification on Compliance Register activities', 'تمكين/تعطيل الإشعارات لأنشطة سجل الامتثال'),
    (@module_id, 'compliance_requirement', 'Enable/Disable notification on Compliance Requirement activities', 'تمكين/تعطيل الإشعارات لأنشطة متطلبات الامتثال'),
    (@module_id, 'compliance_action_plan', 'Enable/Disable notification on Compliance Register Action Plan Activities', 'تمكين/تعطيل الإشعارات لأنشطة خطة عمل سجل الامتثال'),
    (@module_id, 'compliance_requirement_action_plan', 'Enable/Disable notification on Compliance Requirement Action Plan Activities', 'تمكين/تعطيل الإشعارات لأنشطة خطة عمل متطلبات الامتثال');

-- SLA Register notifications
INSERT INTO notification_masters (module_id, identifier, title_en, title_ar)
VALUES
    (@module_id, 'sla_register', 'Enable/Disable notification on SLA Register activities', 'تمكين/تعطيل الإشعارات لأنشطة سجل اتفاقية مستوى الخدمة'),
    (@module_id, 'sla_register_action_plan', 'Enable/Disable notification on SLA Register Action Plan Activities', 'تمكين/تعطيل الإشعارات لأنشطة خطة عمل سجل اتفاقية مستوى الخدمة');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE notification_masters;
-- +goose StatementEnd
