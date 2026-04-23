-- +goose Up
-- +goose StatementBegin
DECLARE @riskModuleID INT;

SET @riskModuleID = (SELECT id FROM modules WHERE identifier = 'risk');

INSERT INTO workflow_contexts (module_id, title_en, title_ar, type) VALUES
(@riskModuleID, 'KRI Register Workflow', 'سير عمل سجل مؤشرات المخاطر الرئيسية', 'kri_register');
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DELETE FROM workflow_contexts WHERE type = 'kri_register';
-- +goose StatementEnd