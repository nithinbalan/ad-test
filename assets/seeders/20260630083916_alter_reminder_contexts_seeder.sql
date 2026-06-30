-- +goose Up
-- +goose StatementBegin

UPDATE reminder_contexts
SET
    type = 'violation-report-action-plan-reminder',
    title_en = 'Violation Report Action Plan Reminder',
    title_ar = N'تذكير خطة عمل تقرير المخالفة',
    module_id = (
        SELECT id
        FROM modules
        WHERE identifier = 'violation_report'
    )
WHERE type = 'awareness-reporting-action-plan-reminder';
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
UPDATE reminder_contexts
SET
    type = 'awareness-reporting-action-plan-reminder',
    title_en = 'Awareness Reporting Action Plan Reminder',
    title_ar = N'تذكير خطة عمل تقارير التوعية',
    module_id = (
        SELECT id
        FROM modules
        WHERE identifier = 'awareness_reporting'
    )
WHERE type = 'violation-report-action-plan-reminder';
-- +goose StatementEnd