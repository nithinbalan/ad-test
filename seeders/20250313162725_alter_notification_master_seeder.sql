-- +goose Up
-- +goose StatementBegin
-- Get the ID of the module with identifier 'business_continuity'
DECLARE @module_id INT;
SELECT @module_id = id FROM modules WHERE identifier = 'business_continuity';

-- Insert notification master records for Compliance Assessment
INSERT INTO notification_masters (module_id, identifier, title_en, title_ar)
VALUES
    (@module_id, 'test_and_exercise', 'Notification on Test And Exercise', 'إشعار بشأن الاختبار والتمرين');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Delete the inserted notification to rollback
DELETE FROM notification_masters WHERE identifier = 'test_and_exercise';
-- +goose StatementEnd
