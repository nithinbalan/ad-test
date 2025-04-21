-- +goose Up
-- +goose StatementBegin
-- Get the ID of the module with identifier 'organizational_context'
DECLARE @module_id INT;
SELECT @module_id = id FROM modules WHERE identifier = 'organizational_context';

-- Insert notification master records for Compliance Assessment
INSERT INTO notification_masters (module_id, identifier, title_en, title_ar)
VALUES
    (@module_id, 'stakeholder_action_plan', 'Notification on Stakeholder Action Plan', 'إشعار بشأن خطة عمل أصحاب المصلحة');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Delete the inserted notification to rollback
DELETE FROM notification_masters WHERE identifier = 'stakeholder_action_plan';
-- +goose StatementEnd
