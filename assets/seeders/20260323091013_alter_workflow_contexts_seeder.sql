-- +goose Up
-- +goose StatementBegin
DECLARE @organizationModuleID INT;

SET @organizationModuleID = (SELECT id FROM modules WHERE identifier = 'organization');

INSERT INTO workflow_contexts (module_id, title_en, title_ar, type) VALUES
(@organizationModuleID, 'Committee Management Workflow', 'سير عمل إدارة اللجنة', 'committee_management');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM workflow_contexts WHERE type = 'committee_management';
-- +goose StatementEnd
