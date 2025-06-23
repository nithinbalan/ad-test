-- +goose Up
-- +goose StatementBegin
INSERT INTO timeline_statuses
    (title_en, title_ar, color_code, type)
VALUES
    ('Pending', 'قيد الانتظار', '#FFC107', 'pending'),
    ('Completed', 'مكتمل', '#28A745', 'completed');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE timeline_statuses;
-- +goose StatementEnd

