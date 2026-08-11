-- +goose Up
-- +goose StatementBegin

DECLARE @module_id INT;
SET @module_id = (
    SELECT id
    FROM modules
    WHERE identifier = 'identity_access_review'
);

IF NOT EXISTS (
    SELECT 1
    FROM notification_masters
    WHERE identifier = 'identity_access_review_campaign'
)
BEGIN
    INSERT INTO notification_masters (
        module_id,
        identifier,
        title_en,
        title_ar
    )
    VALUES (
        @module_id,
        'identity_access_review_campaign',
        'Notification on Identity Access Review Campaign Activities',
        N'إشعار بشأن أنشطة حملة مراجعة الهوية والوصول'
    );
END;

IF NOT EXISTS (
    SELECT 1
    FROM notification_masters
    WHERE identifier = 'identity_access_review_action_plan'
)
BEGIN
    INSERT INTO notification_masters (
        module_id,
        identifier,
        title_en,
        title_ar
    )
    VALUES (
        @module_id,
        'identity_access_review_action_plan',
        'Notification on Identity Access Review Action Plan Activities',
        N'إشعار بشأن أنشطة خطة عمل مراجعة الهوية والوصول'
    );
END;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DELETE FROM notification_masters
WHERE identifier IN (
    'identity_access_review_campaign',
    'identity_access_review_action_plan'
);

-- +goose StatementEnd
