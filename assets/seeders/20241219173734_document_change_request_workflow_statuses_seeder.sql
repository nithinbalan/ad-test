-- +goose Up
-- +goose StatementBegin
INSERT INTO document_change_request_workflow_statuses (title_en, title_ar, color_code, type, status)
VALUES
('Commented', 'حالة التعليق', '#17A2B8', 'comment', 'active'),
('Reverted', 'حالة العودة', '#FF9F45', 'revert', 'active'),
('Rejected', 'حالة الرفض', '#DC3545', 'reject', 'active'),
('Approved', 'حالة الموافقة', '#90DD75', 'approve', 'active'),
('Pending', 'قيد الانتظار', '#FFC107', 'pending', 'active'),
('Resubmitted', 'حالة إعادة التقديم', '#28A745', 'resubmit', 'active');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE document_change_request_workflow_statuses;
-- +goose StatementEnd
