-- +goose Up
-- +goose StatementBegin

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

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

-- Delete re_open record
DELETE FROM workflow_statuses
WHERE type = 'reopen';

-- +goose StatementEnd