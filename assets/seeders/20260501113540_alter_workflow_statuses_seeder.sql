-- +goose Up
-- +goose StatementBegin

-- Drop constraint first (to avoid check constraint issue)
ALTER TABLE workflow_statuses
DROP CONSTRAINT CHK_workflow_statuses_type;

-- Update existing record if already seeded
UPDATE workflow_statuses
SET type = 'reopen',
    title_en = 'ReOpened',
    title_ar = N'إعادة فتح'
WHERE type = 'changes_requested';

-- Insert if not exists (fresh db case)
IF NOT EXISTS (SELECT 1 FROM workflow_statuses WHERE type = 'reopen')
BEGIN
    INSERT INTO workflow_statuses (title_en, title_ar, color_code, type, status)
    VALUES ('ReOpened', N'إعادة فتح', '#FF9F45', 'reopen', 'active');
END

-- Recreate constraint with re_open
ALTER TABLE workflow_statuses
ADD CONSTRAINT CHK_workflow_statuses_type
CHECK (
    type IN (
        'submit',
        'revert',
        'reject',
        'approve',
        'publish',
        'pending',
        'comment',
        'reopen'
    )
);

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

ALTER TABLE workflow_statuses
DROP CONSTRAINT CHK_workflow_statuses_type;

-- Delete re_open record
DELETE FROM workflow_statuses
WHERE type = 'reopen';

-- Restore constraint with 
ALTER TABLE workflow_statuses
ADD CONSTRAINT CHK_workflow_statuses_type
CHECK (
    type IN (
        'submit',
        'revert',
        'reject',
        'approve',
        'publish',
        'pending',
        'comment'
    )
);

-- +goose StatementEnd