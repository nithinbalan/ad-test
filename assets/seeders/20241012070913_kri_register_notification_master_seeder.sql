-- +goose Up
-- +goose StatementBegin
DECLARE @module_id INT;
SET @module_id = (SELECT id FROM modules WHERE identifier = 'risk');

-- Risk module notifications
INSERT INTO notification_masters (module_id, identifier, title_en, title_ar)
VALUES
    (@module_id, 'kri_register', 'Notification on KRI Register', 'الإخطار على سجل KRI');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE notification_masters;
-- +goose StatementEnd
