-- +goose Up
-- +goose StatementBegin
DECLARE @kpiModuleID INT;

SELECT @kpiModuleID = id
FROM modules
WHERE identifier = 'key_performance_indicator';

INSERT INTO reminders (
    module_id,
    title_en,
    title_ar,
    type,
    status
)
VALUES (
    @kpiModuleID,
    'User Tasks',
    N'مهام المستخدم',
    'user-task',
    'active'
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM reminders
WHERE module_id = (SELECT id FROM modules WHERE identifier = 'key_performance_indicator')
  AND type = 'user-task';
-- +goose StatementEnd
