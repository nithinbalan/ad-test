-- +goose Up
-- +goose StatementBegin
DECLARE @ViolationModuleID INT;

SELECT @ViolationModuleID = id
FROM modules
WHERE identifier = 'violation_report';

INSERT INTO reminders (
    module_id,
    title_en,
    title_ar,
    type,
    status
)
VALUES (
      @ViolationModuleID,
       'User Tasks',
      N'مهام المستخدم',
       'user-task',
       'active'
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE r
FROM reminders r
JOIN modules m ON m.id = r.module_id
WHERE m.identifier = 'violation_report'
  AND r.type = 'user-task';

-- +goose StatementEnd
