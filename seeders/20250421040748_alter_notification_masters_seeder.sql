-- +goose Up
-- +goose StatementBegin
-- Get the ID of the module with identifier 'strategy'
DECLARE @module_id INT;
SELECT @module_id = id FROM modules WHERE identifier = 'strategy';

-- Insert notification master records for Action Plan
INSERT INTO notification_masters (module_id, identifier, title_en, title_ar)
VALUES
    (@module_id, 'strategy_profile_action_plan', 'Notification on Strategy Profile Action Plan', 'إشعار بشأن خطة عمل ملف الاستراتيجية');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Delete the inserted notification to rollback
DELETE FROM notification_masters WHERE identifier = 'strategy_profile_action_plan';
-- +goose StatementEnd
