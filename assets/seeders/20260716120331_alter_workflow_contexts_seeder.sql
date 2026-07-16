-- +goose Up
-- +goose StatementBegin
DECLARE @assetModuleID INT;

SET @assetModuleID = (
    SELECT id
    FROM modules
    WHERE identifier = 'asset'
);

INSERT INTO workflow_contexts (
    module_id,
    title_en,
    title_ar,
    type
)
VALUES (
    @assetModuleID,
    'Vulnerability Reject Workflow',
    N'سير عمل رفض الثغرات',
    'vulnerability-reject'
);
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DELETE FROM workflow_contexts
WHERE type = 'vulnerability-reject';
-- +goose StatementEnd