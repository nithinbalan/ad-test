-- +goose Up
-- +goose StatementBegin

DECLARE @module_id INT;
SET @module_id = (SELECT id FROM modules WHERE identifier = 'asset');

INSERT INTO notification_masters (
    module_id,
    identifier,
    title_en,
    title_ar
)
VALUES (
    @module_id,
    'asset_risk_action_plan',
    'Notification on Asset Risk Action Plan Activities',
    N'إشعار بشأن أنشطة خطة عمل مخاطر الأصول'
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DELETE FROM notification_masters
WHERE identifier = 'asset_risk_action_plan';

-- +goose StatementEnd