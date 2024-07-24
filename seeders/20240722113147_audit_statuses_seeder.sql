-- +goose Up
-- +goose StatementBegin
INSERT INTO audit_statuses
    (title_en, title_ar, type, color_code, status)
VALUES
    ('Pending', 'قيد الانتظار', 'pending', '#DC3545', 'active'),
    ('Planned', 'مخطط', 'planned', '#FFA500', 'active'),
    ('In Progress', 'قيد التنفيذ', 'in-progress', '#FFC107', 'active'),
    ('Completed', 'مكتمل', 'completed', '#28A745', 'active');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE audit_statuses;
-- +goose StatementEnd
