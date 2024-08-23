-- +goose Up
-- +goose StatementBegin
INSERT INTO meeting_statuses
    (title_en, title_ar, color_code, type, status)
VALUES
    ('Scheduled', 'مجدول', '#FFA500', 'scheduled', 'active'),
    ('Completed', 'مكتمل', '#28A745', 'completed', 'active'),
    ('Cancelled', 'ملغى', '#DC3545', 'cancelled', 'active');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE meeting_statuses;
-- +goose StatementEnd
