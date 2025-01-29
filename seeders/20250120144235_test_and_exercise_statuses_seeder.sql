-- +goose Up
-- +goose StatementBegin
INSERT INTO test_and_exercise_statuses (title_en, title_ar, color_code, type, status) VALUES
('Scheduled', 'المقرر', '#90DD75', 'scheduled', 'active'),
('Completed', 'مكتمل', '#6C757D', 'completed', 'active');
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE test_and_exercise_statuses;
-- +goose StatementEnd