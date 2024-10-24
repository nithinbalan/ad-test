-- +goose Up
-- +goose StatementBegin
INSERT INTO document_change_request_statuses (title_en, title_ar, color_code, type, status) VALUES
('Initiated', 'مبادر', '#17A2B8', 'initiated', 'active'),   
('Rejected', 'مرفوض', '#E0554D', 'rejected', 'active'),      
('Implementing', 'تنفيذ', '#FD7E14', 'implementing', 'active'),  
('Completed', 'مكتمل', '#90DD75', 'completed', 'active'),    
('Cancelled', 'ملغى', '#6C757D', 'cancelled', 'active');    
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM dbo.document_change_request_statuses WHERE type IN ('initiated', 'rejected', 'implementing', 'completed', 'cancelled');
-- +goose StatementEnd
