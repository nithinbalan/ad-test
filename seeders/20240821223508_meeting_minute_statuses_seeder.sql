-- +goose Up
-- +goose StatementBegin
INSERT INTO meeting_minute_statuses (title_en, title_ar, type, color_code, status) VALUES
('Scheduled', 'مجدول', 'scheduled', '#FFA500', 'active'),
('Cancelled', 'ملغي', 'cancelled', '#DC3545', 'active'),
('Completed', 'مكتمل', 'completed', '#28A745', 'active');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE meeting_minute_statuses;
-- +goose StatementEnd
