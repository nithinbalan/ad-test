
-- +goose Up
-- +goose StatementBegin
INSERT INTO document_workflow_statuses (title_en, title_ar, color_code, type, status)
VALUES
('Commented', 'حالة التعليق', '#D6A2E8', 'comment', 'active'),  
('Reverted', 'حالة العودة', '#FFA500', 'revert', 'active'),    
('Rejected', 'حالة الرفض', '#E0554D', 'reject', 'active'),     
('Approved', 'حالة الموافقة', '#28A745', 'approve', 'active'),
('Pending', 'قيد الانتظار', '#FFC107', 'pending', 'active'); 
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE document_workflow_statuses;
-- +goose StatementEnd
