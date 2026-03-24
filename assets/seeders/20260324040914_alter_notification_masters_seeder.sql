-- +goose Up
-- +goose StatementBegin

-- Get the ID of the module with identifier 'incident'
DECLARE @module_id INT;
SET @module_id = (SELECT id FROM modules WHERE identifier = 'organization');

-- Incident module notifications
INSERT INTO notification_masters (module_id, identifier, title_en, title_ar)
VALUES 
    (@module_id, 'committee_management', 'Notification on Committee Management Activities', 'إشعار بشأن أنشطة إدارة اللجنة');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM notification_masters WHERE identifier = 'committee_management';
-- +goose StatementEnd