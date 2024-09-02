-- +goose Up
-- +goose StatementBegin
INSERT INTO statuses
    (title_en, title_ar, type, color_code, status)
VALUES
    ('Pending', 'قيد الانتظار', 'pending', '#DC3545', 'active'),
    ('Completed', 'مكتمل', 'completed', '#28A745', 'active');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE statuses;
-- +goose StatementEnd
