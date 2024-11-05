-- +goose Up
-- +goose StatementBegin
INSERT INTO threat_statuses (title_en, title_ar, color_code, type, status) VALUES
('New', 'جديد', '#90DD75', 'new', 'active'),
('In Progress', 'قيد التنفيذ', '#FFC107', 'in-progress', 'active'),
('Resolved', 'تم الحل', '#007BFF', 'resolved', 'active'),
('Closed', 'مغلق', '#DC3545', 'closed', 'active');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE threat_statuses;
-- +goose StatementEnd
