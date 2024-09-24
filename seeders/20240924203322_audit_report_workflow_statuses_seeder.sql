
-- +goose Up
-- +goose StatementBegin
INSERT INTO audit_workflow_statuses (title_en, title_ar, color_code, type, status)
VALUES
('Reverted', 'حالة العودة', '#FF9F45', 'revert', 'active'),    
('Rejected', 'حالة الرفض', '#DC3545', 'reject', 'active'),     
('Approved', 'حالة الموافقة', '#90DD75', 'approve', 'active'),
('Pending', 'قيد الانتظار', '#FFC107', 'pending', 'active'); 
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE audit_workflow_statuses;
-- +goose StatementEnd
