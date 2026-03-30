-- +goose Up
-- +goose StatementBegin
INSERT INTO strategic_initiative_statuses
    (title_en, title_ar, color_code, type, status)
VALUES
    ('Planned',     N'مخطط',        '#6C757D', 'planned',     'active'),
    ('In Progress', N'قيد التنفيذ', '#007BFF', 'in_progress', 'active'),
    ('Completed',   N'مكتمل',      '#28A745', 'completed',   'active'),
    ('On Hold',     N'معلق',       '#FFC107', 'on_hold',     'active');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE strategic_initiative_statuses;
-- +goose StatementEnd