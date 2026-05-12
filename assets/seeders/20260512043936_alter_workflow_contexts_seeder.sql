-- +goose Up
-- +goose StatementBegin
DECLARE @riskModuleID INT;

SET @riskModuleID = (SELECT id FROM modules WHERE identifier = 'risk');

INSERT INTO workflow_contexts (module_id, title_en, title_ar, type) VALUES
(@riskModuleID, 'Risk Register Update Plan Workflow', 'سير عمل تحديث خطة المخاطر', 'risk_register_update_plan');
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DELETE FROM workflow_contexts WHERE type = 'risk_register_update_plan';
-- +goose StatementEnd