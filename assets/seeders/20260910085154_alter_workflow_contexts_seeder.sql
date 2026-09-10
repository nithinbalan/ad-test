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
    'Project Security Gate Exception Workflow', 
    N'سير عمل استثناء بوابة أمان المشروع', 
    'project_security_gate_exception'
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM workflow_contexts WHERE type = 'project_security_gate_exception';
-- +goose StatementEnd
