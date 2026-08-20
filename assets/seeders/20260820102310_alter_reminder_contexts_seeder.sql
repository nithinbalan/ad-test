-- +goose Up
-- +goose StatementBegin
DECLARE @kpiModuleID INT;

SELECT @kpiModuleID = id
FROM modules
WHERE identifier = 'key_performance_indicator';

INSERT INTO reminder_contexts (
    type,
    title_en,
    title_ar,
    status,
    module_id
)
VALUES (
    'kpi-review-due-reminder',
    'KPI Review Due Reminder',
    N'استحقاق مراجعة مؤشر الأداء',
    'active',
    @kpiModuleID
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE rr
FROM reminder_rules rr
INNER JOIN reminder_contexts rc ON rc.id = rr.reminder_context_id
WHERE rc.type = 'kpi-review-due-reminder';

DELETE FROM reminder_contexts
WHERE type = 'kpi-review-due-reminder';
-- +goose StatementEnd
