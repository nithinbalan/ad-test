-- +goose Up
-- +goose StatementBegin
DECLARE @riskModuleID INT;

SET @riskModuleID = (SELECT id FROM modules WHERE identifier = 'risk');

INSERT INTO workflow_contexts (module_id, title_en, title_ar, type) VALUES
(@riskModuleID, 'Risk Response Workflow', 'سير عمل استجابة المخاطر', 'risk_response');
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DELETE FROM workflow_contexts WHERE type = 'risk_response';
-- +goose StatementEnd