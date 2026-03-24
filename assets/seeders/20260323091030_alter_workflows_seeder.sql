-- +goose Up
-- +goose StatementBegin
DECLARE @committeeWorkflowContextID INT;

SET @committeeWorkflowContextID = (SELECT id FROM workflow_contexts WHERE type = 'committee_management');

INSERT INTO workflows (status, created_by_id, created_at, title_en, title_ar,workflow_context_id)
VALUES 
    ('active', 1, GETDATE(), 'Committee Management Workflow', 'سير عمل إدارة اللجنة', @committeeWorkflowContextID);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM workflows 
WHERE workflow_context_id = (
    SELECT id FROM workflow_contexts WHERE type = 'committee_management'
);
-- +goose StatementEnd
