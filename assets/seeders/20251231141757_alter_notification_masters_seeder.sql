-- +goose Up
-- +goose StatementBegin

-- Get the ID of the module with identifier 'tprm'
DECLARE @module_id INT;
SET @module_id = (SELECT id FROM modules WHERE identifier = 'tprm');

-- tprm module notifications
INSERT INTO notification_masters (module_id, identifier, title_en, title_ar)
VALUES
    (@module_id, 'tprm_assessment', 'Notification on TPRM Assessment Activities', 'إشعار بشأن أنشطة تقييم إدارة مخاطر الطرف الثالث');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM notification_masters WHERE identifier = 'tprm_assessment';
-- +goose StatementEnd