-- +goose Up
-- +goose StatementBegin
UPDATE workflow_contexts
SET
    type = 'vulnerability-closure',
    title_en = 'Vulnerability Closure Workflow',
    title_ar = N'سير عمل إغلاق الثغرات'
WHERE type = 'vulnerability-reject';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
UPDATE workflow_contexts
SET
    type = 'vulnerability-reject',
    title_en = 'Vulnerability Reject Workflow',
    title_ar = N'سير عمل رفض الثغرات'
WHERE type = 'vulnerability-closure';
-- +goose StatementEnd
