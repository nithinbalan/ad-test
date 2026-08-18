-- +goose Up
-- +goose StatementBegin
DECLARE @incidentModuleID INT;

SET @incidentModuleID = (SELECT id FROM modules WHERE identifier = 'incident');

INSERT INTO workflow_contexts (module_id, title_en, title_ar, type) VALUES
(@incidentModuleID, 'Incident Closure Review Workflow', 'سير عمل مراجعة إغلاق الحادث', 'incident_closure');
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DELETE FROM workflow_contexts WHERE type = 'incident_closure';
-- +goose StatementEnd