-- +goose Up
-- +goose StatementBegin
INSERT INTO third_party_assessment_statuses
    (title_en, title_ar, color_code, type, status)
VALUES
    ('Draft', 'نشيط', '#007BFF', 'draft', 'active'),
    ('New', 'غير نشط', '#FFC107', 'new', 'active'),
    ('In Progress', 'مؤرشف', '#FD7E14', 'in-progress', 'active'),
    ('Submitted', 'مؤرشف', '#28A745', 'submitted', 'active');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE third_party_assessment_statuses;
-- +goose StatementEnd
