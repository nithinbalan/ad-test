-- +goose Up
-- +goose StatementBegin
-- Insert the new "Resubmit" status
INSERT INTO document_workflow_statuses (title_en, title_ar, color_code, type, status)
VALUES
('Resubmitted', 'حالة إعادة التقديم', '#28A745', 'resubmit', 'active'); 
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Remove the "Resubmit" status
DELETE FROM document_workflow_statuses
WHERE type = 'resubmit';
-- +goose StatementEnd
