-- +goose Up
-- +goose StatementBegin
DECLARE @AssetModuleID INT;

SELECT @AssetModuleID = id
FROM modules
WHERE identifier = 'asset';

INSERT INTO reminder_contexts (
    type,
    title_en,
    title_ar,
    status,
    module_id
)
VALUES (
    'vulnerability-action-plan-reminder',
    'Vulnerability Action Plan Reminder',
    N'تذكير بخطة عمل الثغرات',
    'active',
    @AssetModuleID
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM reminder_contexts
WHERE type = 'vulnerability-action-plan-reminder';
-- +goose StatementEnd
