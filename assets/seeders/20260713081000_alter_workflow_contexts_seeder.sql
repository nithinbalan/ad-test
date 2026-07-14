-- +goose Up
-- +goose StatementBegin
DECLARE @assetModuleID INT;

SET @assetModuleID = (SELECT id FROM modules WHERE identifier = 'asset');

INSERT INTO workflow_contexts (module_id, title_en, title_ar, type) VALUES
(@assetModuleID, 'Risk Acceptance Workflow', 'سير عمل قبول المخاطر', 'asset_risk');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM workflow_contexts WHERE type = 'asset_risk';
-- +goose StatementEnd
