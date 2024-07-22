-- +goose Up
-- +goose StatementBegin
INSERT INTO audit_statuses
    (title_en, title_ar, type, color_code, status)
VALUES
    ('Pending', 'قيد الانتظار', 'pending', '#FFFF00', 'active'),
    ('Planned', 'مخطط', 'planned', '#FFA500', 'active'),
    ('In Progress', 'قيد التنفيذ', 'in-progress', '#0000FF', 'active'),
    ('Completed', 'مكتمل', 'completed', '#00FF00', 'active');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE audit_statuses;
-- +goose StatementEnd
