-- +goose Up
-- +goose StatementBegin
DECLARE @complianceModuleID INT;

SET @complianceModuleID = (SELECT id FROM modules WHERE identifier = 'compliance');

INSERT INTO workflow_contexts (module_id, title_en, title_ar, type) VALUES
(@complianceModuleID, 'Cyber Compliance Register Workflow', 'سير عمل سجل الامتثال السيبراني', 'cyber_compliance_register');
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DELETE FROM workflow_contexts WHERE type = 'cyber_compliance_register';
-- +goose StatementEnd