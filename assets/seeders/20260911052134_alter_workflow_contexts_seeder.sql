-- +goose Up
-- +goose StatementBegin
DECLARE @projectModuleID INT;

SET @projectModuleID = (SELECT id FROM modules WHERE identifier = 'project');

INSERT INTO workflow_contexts (
    module_id, 
    title_en, 
    title_ar, 
    type
)
VALUES (
    @projectModuleID, 
    'Project Register Workflow', 
    N'سير عمل سجل المشاريع', 
    'project_register'
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM workflow_contexts WHERE type = 'project_register';
-- +goose StatementEnd
