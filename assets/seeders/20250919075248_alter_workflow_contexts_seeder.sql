-- +goose Up
-- +goose StatementBegin
DECLARE @riskModuleID INT;

SET @riskModuleID = (SELECT id FROM modules WHERE identifier = 'risk');

INSERT INTO workflow_contexts (module_id, title_en, title_ar, type) VALUES
(@riskModuleID, 'Risk Workflow', 'سير عمل المخاطر', 'risk_register');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM workflow_contexts WHERE type = 'risk_register';
-- +goose StatementEnd
