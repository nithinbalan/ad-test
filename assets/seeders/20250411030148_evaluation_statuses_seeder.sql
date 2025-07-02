-- +goose Up
-- +goose StatementBegin
INSERT INTO evaluation_statuses
    (title_en, title_ar, color_code, type, status)
VALUES
    ('Fully Met', 'التقى', '#90DD75', 'fully-met', 'active'),
    ('Partially Met', 'تم الوفاء بها جزئيًا', '#FF9F45', 'partially-met', 'active'),
    ('Not Met', 'لم يلتق', '#007BFF', 'not-met', 'active');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE evaluation_statuses;
-- +goose StatementEnd
