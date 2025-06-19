-- +goose Up
-- +goose StatementBegin
INSERT INTO competency_levels (title_en, title_ar, status) VALUES 
('Beginner', 'مبتدئ', 'active'),
('Intermediate', 'متوسط', 'active'),
('Advanced', 'متقدم', 'active');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM competency_levels WHERE title_en IN ('Beginner', 'Intermediate', 'Advanced');
-- +goose StatementEnd
