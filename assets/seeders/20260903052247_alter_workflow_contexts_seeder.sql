-- +goose Up
-- +goose StatementBegin
DECLARE @projectModuleID INT;

SET @projectModuleID = (SELECT id FROM modules WHERE identifier = 'project');

INSERT INTO workflow_contexts (module_id, title_en, title_ar, type) VALUES
(@projectModuleID, 'Project Change Request Workflow', N'سير عمل طلب تغيير المشروع', 'project_change_request');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM workflow_contexts WHERE type = 'project_change_request';
-- +goose StatementEnd
