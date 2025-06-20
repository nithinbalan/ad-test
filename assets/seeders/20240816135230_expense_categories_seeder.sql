-- +goose Up
-- +goose StatementBegin
INSERT INTO expense_categories
    (title_en, title_ar, description_en, description_ar, status)
VALUES
    ('Software', 'برمجيات', 'Software-related expenses', 'نفقات متعلقة بالبرمجيات', 'active'),
    ('Hardware', 'أجهزة', 'Hardware-related expenses', 'نفقات متعلقة بالأجهزة', 'active'),
    ('Training', 'تدريب', 'Training-related expenses', 'نفقات متعلقة بالتدريب', 'active'),
    ('Consultancy', 'استشارات', 'Consultancy-related expenses', 'نفقات متعلقة بالاستشارات', 'active');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE expense_categories;
-- +goose StatementEnd
