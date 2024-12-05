-- +goose Up
-- +goose StatementBegin
-- Insert the comment status into the workflow_statuses table
INSERT INTO workflow_statuses (title_en, title_ar, color_code, type, status)
VALUES
('Comment', 'تعليق', '#17A2B8', 'comment', 'active');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Remove the comment status from the workflow_statuses table
DELETE FROM workflow_statuses
WHERE type = 'comment';
-- +goose StatementEnd
