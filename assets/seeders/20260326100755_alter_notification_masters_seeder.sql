-- +goose Up
-- +goose StatementBegin

-- Get the ID of the module with identifier 'consultation_request'
DECLARE @module_id INT;
SET @module_id = (SELECT id FROM modules WHERE identifier = 'consultation_request');

-- consultation request module notifications
INSERT INTO notification_masters (module_id, identifier, title_en, title_ar)
VALUES
    (@module_id, 'consultation_request_register', 'Notification on Consultation Request activities', 'إشعار بشأن أنشطة طلبات الاستشارة');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM notification_masters WHERE identifier = 'consultation_request_register';
-- +goose StatementEnd