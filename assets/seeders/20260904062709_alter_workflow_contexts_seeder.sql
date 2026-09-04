-- +goose Up
-- +goose StatementBegin
DECLARE @projectModuleID INT;

SET @projectModuleID = (SELECT id FROM modules WHERE identifier = 'project');

INSERT INTO workflow_contexts (module_id, title_en, title_ar, type) VALUES
(@projectModuleID, 'Project Decision Workflow', N'سير عمل قرارات المشروع', 'project_decision_log');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM workflow_contexts WHERE type = 'project_decision_log';
-- +goose StatementEnd
