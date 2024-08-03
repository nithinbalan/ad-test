
-- +goose Up
-- +goose StatementBegin
INSERT INTO document_workflow_statuses (title_en, title_ar, color_code, type, status)
VALUES
('Comment', 'حالة التعليق', '#D6A2E8', 'comment', 'active'),  
('Revert', 'حالة العودة', '#FFA500', 'revert', 'active'),    
('Reject', 'حالة الرفض', '#E0554D', 'reject', 'active'),     
('Approve', 'حالة الموافقة', '#28A745', 'approve', 'active'),
('Pending', 'قيد الانتظار', '#FFC107', 'pending', 'active'); 
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE document_workflow_statuses;
-- +goose StatementEnd
