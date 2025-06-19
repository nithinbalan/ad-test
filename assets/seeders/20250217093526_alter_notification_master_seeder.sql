-- +goose Up
-- +goose StatementBegin
-- Get the ID of the module with identifier 'audit'
DECLARE @module_id INT;
SET @module_id = (SELECT id FROM modules WHERE identifier = 'compliance');

-- Risk module notifications
INSERT INTO notification_masters (module_id, identifier, title_en, title_ar)
VALUES
    (@module_id, 'compliance_assessment_requirement', 'Notification on Compliance Assessment Requirement Activities', 'إشعار بشأن أنشطة متطلبات تقييم الامتثال');
    
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE notification_masters;
-- +goose StatementEnd
