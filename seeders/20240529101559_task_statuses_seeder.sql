-- +goose Up
-- +goose StatementBegin
INSERT INTO task_statuses
    (title_en, title_ar, color_code, type, is_open, status)
VALUES
    ('New', 'جديد', '#FF0000', 'new', 1, 'active'),
    ('In Progress', 'قيد التنفيذ', '#FFFF00', 'in_progress', 1, 'active'),
    ('Closed', 'مغلق', '#008000', 'closed', 0, 'active');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE task_statuses;
-- +goose StatementEnd
