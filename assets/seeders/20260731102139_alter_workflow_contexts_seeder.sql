-- +goose Up
-- +goose StatementBegin
DECLARE @businessContinuityModuleID INT;

SET @businessContinuityModuleID = (SELECT id FROM modules WHERE identifier = 'business_continuity');

INSERT INTO workflow_contexts (module_id, title_en, title_ar, type) VALUES
(@businessContinuityModuleID, 'Cyber BIA Workflow', 'سير عمل تقييم الأثر على الأعمال السيبراني', 'cyber_bia');
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DELETE FROM workflow_contexts WHERE type = 'cyber_bia';
-- +goose StatementEnd