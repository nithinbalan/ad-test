-- +goose Up
-- +goose StatementBegin
IF NOT EXISTS (
    SELECT 1
    FROM workflow_statuses
    WHERE type = 'changes_requested'
)
AND EXISTS (
    SELECT 1
    FROM sys.check_constraints
    WHERE parent_object_id = OBJECT_ID(N'dbo.workflow_statuses')
      AND name = 'CHK_workflow_statuses_type'
      AND definition LIKE '%''changes_requested''%'
)
BEGIN
    INSERT INTO workflow_statuses
        (title_en, title_ar, color_code, type, status)
    VALUES
        ('Changes Requested', N'تم طلب التغييرات',
         '#FF9F45', 'changes_requested', 'active');
END;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM workflow_statuses
WHERE type = 'changes_requested';
-- +goose StatementEnd
