-- +goose Up
-- +goose StatementBegin
INSERT INTO grc_report_workflow_statuses (title_en, title_ar, color_code, type, status) VALUES
('Draft', 'مسودة', '#0000FF', 'draft', 'active'),
('In Review', 'قيد المراجعة', '#FFD700', 'in-review', 'active'),
('Reviewed', 'تمت المراجعة', '#FF4500', 'reviewed', 'active'),
('Pending Approval', 'في انتظار الموافقة', '#800080', 'pending-approval', 'active'),
('Approved', 'موافقة', '#008000', 'approved', 'active');
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE grc_report_workflow_statuses;
-- +goose StatementEnd