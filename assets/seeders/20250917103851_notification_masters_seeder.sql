-- +goose Up
-- +goose StatementBegin

-- Get the ID of the module with identifier 'organization'
DECLARE @module_id INT;
SET @module_id = (SELECT id FROM modules WHERE identifier = 'organization');

-- User Delegation notifications
INSERT INTO notification_masters (module_id, identifier, title_en, title_ar)
VALUES
    (@module_id, 'user_management', 'Notification on User Delegation', 'إشعار بشأن تفويض المستخدم');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM notification_masters WHERE identifier = 'user_management';
-- +goose StatementEnd
