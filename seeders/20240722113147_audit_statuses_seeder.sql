-- +goose Up
-- +goose StatementBegin
INSERT INTO audit_statuses
    (title_en, title_ar, type, color_code, status)
VALUES
    ('Draft', 'مسودة', 'draft', '#007BFF', 'active'),
    ('Published', 'منشور', 'published', '#007BFF', 'active'),
    ('In Progress', 'قيد التنفيذ', 'in-progress', '#FFC107', 'active'),
    ('Completed', 'مكتمل', 'completed', '#28A745', 'active');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE audit_statuses;
-- +goose StatementEnd

