-- +goose Up
-- +goose StatementBegin

DECLARE @AwarenessReportingModuleID INT;

SELECT @AwarenessReportingModuleID = id
FROM modules
WHERE identifier = 'awareness_reporting';

INSERT INTO reminder_contexts (
    type,
    title_en,
    title_ar,
    status,
    module_id
)
VALUES (
    'awareness-reporting-action-plan-reminder',
    'Awareness Reporting Action Plan Reminder',
    N'تذكير خطة عمل تقارير التوعية',
    'active',
    @AwarenessReportingModuleID
);

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

DELETE FROM reminder_contexts
WHERE type = 'awareness-reporting-action-plan-reminder';

-- +goose StatementEnd