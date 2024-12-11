-- +goose Up
-- +goose StatementBegin
INSERT INTO document_change_request_statuses (title_en, title_ar, color_code, type, status) VALUES
('Draft', 'مسودة', '#90DD75', 'draft', 'active'),
('In Review', 'قيد المراجعة', '#17A2B8', 'in-review', 'active'),
('Rejected', 'مرفوض', '#E0554D', 'rejected', 'active'),
('Approved', 'موافق عليه', '#FFC107', 'approved', 'active'),
('Published', 'منشور', '#6C757D', 'published', 'active'),
('Archived', 'مؤرشف', '#FD7E14', 'archived', 'active'),
('Reverted', 'معاد', '#28A745', 'reverted', 'active');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE document_change_request_statuses;
-- +goose StatementEnd
