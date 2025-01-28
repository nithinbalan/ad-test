-- +goose Up
-- +goose StatementBegin
INSERT INTO test_and_exercise_statuses (title_en, title_ar, color_code, type, status) VALUES
('In Progress', 'في تَقَدم', '#FFC107', 'in-progress', 'active');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE document_change_request_statuses;
-- +goose StatementEnd
