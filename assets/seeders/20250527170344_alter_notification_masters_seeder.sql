-- +goose Up
-- +goose StatementBegin
-- Get the ID of the module with identifier 'strategy'
DECLARE @module_id INT;
SELECT @module_id = id FROM modules WHERE identifier = 'report';

-- Insert notification master records for Action Plan
INSERT INTO notification_masters (module_id, identifier, title_en, title_ar)
VALUES
    (@module_id, 'grc_report', 'Notification on GRC Report', 'إشعار بشأن تقرير GRC');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Delete the inserted notification to rollback
DELETE FROM notification_masters WHERE identifier = 'grc_report';
-- +goose StatementEnd
