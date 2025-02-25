-- +goose Up
-- +goose StatementBegin
-- Get the ID of the module with identifier 'compliance'
DECLARE @module_id INT;
SELECT @module_id = id FROM modules WHERE identifier = 'compliance';

-- Insert notification master records for Compliance Assessment
INSERT INTO notification_masters (module_id, identifier, title_en, title_ar)
VALUES
    (@module_id, 'compliance_assessment', 'Notification on Compliance Assessment', 'إشعار بشأن تقييم الامتثال');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Delete the inserted notification to rollback
DELETE FROM notification_masters WHERE identifier = 'compliance_assessment';
-- +goose StatementEnd
