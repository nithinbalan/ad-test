-- +goose Up
-- +goose StatementBegin
INSERT INTO awareness_maturity_levels (title_en, title_ar, type)
VALUES
('Excellent', N'ممتاز', 'excellent'),
('Very Good', N'جيد جداً', 'very-good'),
('Good', N'جيد', 'good'),
('Acceptable', N'مقبول', 'acceptable');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM awareness_maturity_levels;
-- +goose StatementEnd
