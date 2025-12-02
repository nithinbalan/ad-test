-- +goose Up
-- +goose StatementBegin
DECLARE @assessmentModuleID INT;

SET @assessmentModuleID = (SELECT id FROM modules WHERE identifier = 'tprm');

INSERT INTO workflow_contexts (module_id, title_en, title_ar, type) VALUES
(@assessmentModuleID, 'Assessment Workflow', 'سير عمل التقييم', 'tprm_assessment');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM workflow_contexts WHERE type = 'tprm_assessment';
-- +goose StatementEnd
