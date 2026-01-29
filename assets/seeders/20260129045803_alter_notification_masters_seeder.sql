-- +goose Up
-- +goose StatementBegin

-- Get the ID of the module with identifier 'violation report'
DECLARE @module_id INT;
SET @module_id = (SELECT id FROM modules WHERE identifier = 'violation_report');

-- violation report module notifications
INSERT INTO notification_masters (module_id, identifier, title_en, title_ar)
VALUES
    (@module_id, 'violation_report', 'Notification on Violation Report activities', 'إشعار بشأن أنشطة تقارير الانتهاك'), 
    (@module_id, 'violation_report_action_plan', 'Notification on Violation Report Action Plan Activities', 'إشعار بشأن أنشطة خطة عمل التقارير');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM notification_masters WHERE identifier = 'violation_report';
DELETE FROM notification_masters WHERE identifier = 'violation_report_action_plan';
-- +goose StatementEnd
