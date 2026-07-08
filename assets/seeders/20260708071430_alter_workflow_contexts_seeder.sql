-- +goose Up
-- +goose StatementBegin
DECLARE @businessContinuityModuleID INT;

SET @businessContinuityModuleID = (SELECT id FROM modules WHERE identifier = 'business_continuity');

INSERT INTO workflow_contexts (module_id, title_en, title_ar, type) VALUES
(@businessContinuityModuleID, 'BIA Workflow', 'سير عمل تقييم الأثر على الأعمال', 'bia');
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DELETE FROM workflow_contexts WHERE type = 'bia';
-- +goose StatementEnd