-- +goose Up
-- +goose StatementBegin
INSERT INTO planned_advisory_importances (title_en, title_ar, status) VALUES
    ('Immediate', N'فوري', 'active'),
    ('Highly Urgent', N'عاجل جدا', 'active'),
    ('Urgent', N'عاجل', 'active'),
    ('Normal', N'عادي', 'active');
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DELETE FROM planned_advisory_importances
WHERE title_en IN ('Immediate', 'Highly Urgent', 'Urgent', 'Normal');
-- +goose StatementEnd
