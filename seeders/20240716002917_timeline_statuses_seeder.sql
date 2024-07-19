-- +goose Up
-- +goose StatementBegin
INSERT INTO timeline_statuses
    (title_en, title_ar, color_code, type)
VALUES
    ('Pending', 'قيد الانتظار', '#FFA500', 'pending'),
    ('Completed', 'مكتمل', '#008000', 'completed');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE timeline_statuses;
-- +goose StatementEnd

