-- +goose Up
-- +goose StatementBegin
INSERT INTO document_statuses (title_en, title_ar, color_code, type, status) VALUES
('Draft', 'مسودة', '#F0F0F0', 'draft', 'active'),
('In Review', 'قيد المراجعة', '#FFFF00', 'in-review', 'active'),
('Rejected', 'مرفوض', '#FF0000', 'rejected', 'active'),
('Approved', 'موافق عليه', '#008000', 'approved', 'active'),
('Published', 'منشور', '#0000FF', 'published', 'active'),
('Archived', 'مؤرشف', '#808080', 'archived', 'active');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE document_statuses;
-- +goose StatementEnd
