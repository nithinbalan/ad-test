-- +goose Up
-- +goose StatementBegin

DECLARE @contextID INT;

SET @contextID = (
    SELECT id
    FROM workflow_contexts
    WHERE [type] = 'asset_risk_response'
);

INSERT INTO workflows (
    type,
    status,
    created_by_id,
    created_at,
    title_en,
    title_ar,
    workflow_context_id
)
VALUES (
    'asset-risk-response-workflow',
    'active',
    1,
    GETDATE(),
    'Asset Risk Response Workflow',
    N'سير عمل معالجة مخاطر الأصول',
    @contextID
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DECLARE @contextID INT;

SET @contextID = (
    SELECT id
    FROM workflow_contexts
    WHERE [type] = 'asset_risk_response'
);

DELETE FROM workflows
WHERE workflow_context_id = @contextID;

-- +goose StatementEnd