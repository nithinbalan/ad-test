-- +goose Up
-- +goose StatementBegin

-- Get the ID of the module with identifier 'document'
DECLARE @module_id INT;
SET @module_id = (SELECT id FROM modules WHERE identifier = 'document');

-- Document module notifications
INSERT INTO notification_masters (module_id, identifier, title_en, title_ar)
VALUES
    (@module_id, 'document_change_request', 'Notification on Document Change Request Activities', 'إشعار بشأن أنشطة طلب تغيير المستند');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM notification_masters WHERE identifier = 'document_change_request';
-- +goose StatementEnd