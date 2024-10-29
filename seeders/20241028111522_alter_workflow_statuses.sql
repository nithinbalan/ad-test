-- +goose Up
-- +goose StatementBegin
-- Insert the pending status into the workflow_statuses table
INSERT INTO workflow_statuses (title_en, title_ar, color_code, type, status)
VALUES
('Pending', 'قيد الانتظار', '#FFC107', 'pending', 'active');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Remove the pending status from the workflow_statuses table
DELETE FROM workflow_statuses
WHERE type = 'pending';
-- +goose StatementEnd
