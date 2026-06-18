-- +goose Up
-- +goose StatementBegin
INSERT INTO awareness_evaluation_ratings
(title_en, title_ar, type, status)
VALUES
    ('Excellent', N'ممتاز', 'excellent', 'active'),
    ('Good', N'جيد', 'good', 'active'),
    ('Average', N'متوسط', 'average', 'active'),
    ('Poor', N'ضعيف', 'poor', 'active');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM awareness_evaluation_ratings;
-- +goose StatementEnd